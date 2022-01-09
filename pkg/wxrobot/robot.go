package wxrobot

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kuangcp/logger"
)

type (
	Content struct {
		Content             string   `json:"content"`
		MentionedList       []string `json:"mentioned_list,omitempty"`        // 仅text类型 有效
		MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"` // 仅text类型 有效
	}

	ArticleList struct {
		Articles []Article `json:"articles"`
	}

	Article struct {
		Title       string `json:"title"`       // 标题，不超过128个字节，超过会自动截断
		Description string `json:"description"` // 描述，不超过512个字节，超过会自动截断
		URL         string `json:"url"`         // 点击后跳转的链接
		PicURL      string `json:"picurl"`      // 支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150
	}

	Image struct {
		Base64 string `json:"base64"` // 图片内容的base64编码
		MD5    string `json:"md5"`    // 图片内容（base64编码前）的md5值
	}

	Msg struct {
		MsgType string `json:"msgtype"`
		// 以下字段为单选值
		Text     *Content     `json:"text,omitempty"`     // content 最长不超过2048个字节，必须是utf8编码
		Markdown *Content     `json:"markdown,omitempty"` // content 最长不超过4096个字节，必须是utf8编码
		News     *ArticleList `json:"news,omitempty"`
		Image    *Image       `json:"image,omitempty"`
	}

	// Robot 文档 https://work.weixin.qq.com/api/doc/90000/90136/91770
	//  1. 接口调用限流：20条消息/min
	//  2. 当前自定义机器人支持 文本（text）、markdown（markdown）、图片（image）、图文（news）四种消息类型。
	//  3. 机器人的text/markdown类型消息支持在content中使用<@userid>扩展语法来@群成员
	Robot struct {
		SecretKey   string
		MockRequest bool
		limiter     *PeriodRateLimiter
		client      *http.Client
	}
)

const (
	robotApi            = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	imgMaxSize          = 2 << 20 // 发送图片最大 2Mib
	imgToBase64SizeRate = 1.34    // 图片转base64 空间膨胀比例为 4/3
)

func NewRobot(key string) *Robot {
	return &Robot{SecretKey: key, limiter: NewCustomLimiter(time.Minute, 19), client: &http.Client{}}
}

// sendJSONPost 发送body为JSON的 Post 请求
//  return response,timeWasted,error
func (r *Robot) sendJSONPost(value interface{}) ([]byte, int64, error) {
	start := time.Now()
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, 0, err
	}

	reader := bytes.NewReader(jsonBytes)
	request, err := http.NewRequest("POST", robotApi+r.SecretKey, reader)
	if err != nil {
		return nil, 0, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	if r.MockRequest {
		logger.Info(start, "\n", string(jsonBytes))
		return nil, 0, nil
	}

	resp, err := r.client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return rspBody, time.Now().Sub(start).Milliseconds(), nil
}

// SendMarkDown markdown消息
//  支持格式：
//   1-6级标题:   #
//   加粗:       **文字**
//   行内代码段:  ``
//   链接:       [文字](URL)
//   引用:       > 文字
//   字体颜色:    <font color="info">绿色</font> <font color="comment">灰色</font> <font color="warning">橙红色</font>
func (r *Robot) SendMarkDown(content Content) error {
	if !r.limiter.acquire() {
		return errors.New("out of limiter")
	}
	msg := Msg{MsgType: "markdown", Markdown: &content}

	result, waste, err := r.sendJSONPost(msg)
	if err != nil {
		return err
	}
	logger.Warn(string(result), " time: ", waste)
	return nil
}

// SendText 文本消息
func (r *Robot) SendText(content Content) error {
	if !r.limiter.acquire() {
		return errors.New("out of limiter")
	}

	msg := Msg{MsgType: "text", Text: &content}

	result, waste, err := r.sendJSONPost(msg)
	if err != nil {
		return err
	}
	logger.Warn(string(result), " time: ", waste)
	return nil
}

// SendNews 发送图文消息
//  注意：单个图文时能看到 title 和 description, 多个时只能看到 title
func (r *Robot) SendNews(articles ...Article) error {
	if !r.limiter.acquire() {
		return errors.New("out of limiter")
	}

	if articles == nil {
		return errors.New("empty param")
	}

	msg := Msg{MsgType: "news", News: &ArticleList{Articles: articles}}
	result, waste, err := r.sendJSONPost(msg)
	if err != nil {
		return err
	}
	logger.Warn(string(result), " time: ", waste)
	return nil
}

// SendImageByFile 发送图片
func (r *Robot) SendImageByFile(filePath string) error {
	if !r.limiter.acquire() {
		return errors.New("out of limiter")
	}

	open, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return r.SendImageByBytes(open)
}

// TODO 优化项 共用缓存
func imgToBase64(img []byte) string {
	dist := make([]byte, int(float32(len(img))*imgToBase64SizeRate)+1)
	base64.StdEncoding.Encode(dist, img) //buff转成base64
	index := bytes.IndexByte(dist, 0)    //这里要注意，由于先申请了固定长度数组，所以没有被填充完的末端需要去掉
	baseImage := dist[0:index]
	return string(baseImage)
}

func buildImgFileMd5(img []byte) string {
	hashFunc := md5.New()
	hashFunc.Write(img)
	return hex.EncodeToString(hashFunc.Sum(nil))
}

// SendImageByBytes 发送图片
//  图片（base64编码前）最大不能超过2M，支持JPG,PNG格式
func (r *Robot) SendImageByBytes(img []byte) error {
	if !r.limiter.acquire() {
		return errors.New("out of limiter")
	}

	if img == nil {
		return errors.New("empty param")
	}
	if len(img) > imgMaxSize {
		return errors.New(fmt.Sprintf("over limit (2Mib): %v byte", len(img)))
	}

	msg := Msg{
		MsgType: "image",
		Image: &Image{
			Base64: imgToBase64(img),
			MD5:    buildImgFileMd5(img),
		},
	}

	result, waste, err := r.sendJSONPost(msg)
	if err != nil {
		return err
	}
	logger.Warn(string(result), " time: ", waste)
	return nil
}

// TODO 文件上传
// TODO 文件消息

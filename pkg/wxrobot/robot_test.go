package wxrobot

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/kuangcp/logger"
)

var appKey = os.Getenv("appkey")

func TestText(t *testing.T) {
	robot := NewRobot(appKey)
	robot.MockRequest()
	robot.SendText(Content{Content: "闹钟", MentionedList: []string{"xx", "@all"}})
}

func TestTextWithLimiter(t *testing.T) {
	robot := NewRobot(appKey)
	robot.MockRequest()

	go loopSend(robot)
	time.Sleep(time.Second * 3)
	go loopSend(robot)
	time.Sleep(time.Second * 6)
	go loopSend(robot)
	wxrobot := robot.(*WeWorkRobot)
	for range time.NewTicker(time.Second).C {
		log.Println(wxrobot.limiter.queueState())
	}
}

func loopSend(robot Robot) {
	for t := range time.NewTicker(time.Second * 4).C {
		time.Sleep(time.Duration(rand.Intn(300)+600) * time.Millisecond)
		err := robot.SendText(Content{Content: "闹钟 " + t.String(), MentionedList: []string{"xx", "@all"}})
		if err != nil {
			log.Println(err)
		}
	}
}

func TestArticle(t *testing.T) {
	robot := NewRobot(appKey)
	logger.Info(robot)
	robot.MockRequest()
	logger.Info(robot)

	umlArti := Article{Title: "UML", Description: "UML 构图", URL: "https://blog.csdn.net/simonezhlx/article/details/8855297", PicURL: "https://img-blog.csdn.net/20130426171100149"}
	blogArti := Article{Title: "Spring", Description: "从设计角度，深入分析 Spring 循环依赖的解决思路", URL: "https://juejin.cn/post/6958989396917895205", PicURL: "https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/5a690e0622054ed1a5463751f300c815~tplv-k3u1fbpfcp-zoom-crop-mark:1304:1304:1304:734.image"}
	robot.SendNews(umlArti, blogArti)
}

func TestImage(t *testing.T) {
	robot := NewRobot(appKey)
	robot.MockRequest()

	//robot.SendImageByFile("/home/kcp/Pictures/2020-06-20_17-19.png")
	//err := robot.SendImageByFile("/home/kcp/Pictures/uml.svg")
	err := robot.SendImageByFile("/home/kcp/Pictures/2020-06-20_11-20.png")
	if err != nil {
		log.Println(err)
	}
}

func TestImageBase64(t *testing.T) {
	dir, err := os.ReadDir("/home/kcp/Pictures/")
	if err != nil {
		return
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		file, err := ioutil.ReadFile("/home/kcp/Pictures/" + entry.Name())
		if err != nil {
			continue
		}
		base64 := imgToBase64(file)
		fmt.Println(len(file), len(base64), float32(len(base64))/float32(len(file)))
	}
}

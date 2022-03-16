package wxrobot

type Robot interface {

	// MockRequest 不触发真实请求
	MockRequest()

	// ShowRequestLog 打印请求参数
	ShowRequestLog()

	MarkDownGrey(content string) string
	MarkDownGreen(content string) string
	MarkDownOrange(content string) string

	// SendMarkDown markdown消息
	//  支持格式：
	//   1-6级标题:   # 文字
	//   加粗:       **文字**
	//   行内代码段:  ``
	//   链接:       [文字](URL)
	//   引用:       > 文字
	//   字体颜色:    <font color="info">绿色</font> <font color="comment">灰色</font> <font color="warning">橙红色</font>
	SendMarkDown(content string) error

	// SendText 文本消息
	// 支持 @ 功能
	SendText(content Content) error

	// SendNews 发送图文消息
	//  注意：单个图文时能看到 title 和 description, 多个时只能看到 title
	SendNews(articles ...Article) error

	// SendImageByFile 发送图片
	// filePath 图片绝对路径 或者 运行目录的相对路径
	SendImageByFile(filePath string) error

	// SendImageByBytes 发送图片
	//  img 图片文件字节（base64编码前）最大不能超过2M，支持JPG,PNG格式
	SendImageByBytes(img []byte) error
}

package wxrobot

const (
	MdColorGreen  = "info"
	MdColorGray   = "comment"
	MdColorOrange = "warning"

	imgMaxSize          = 2 << 20 // 发送图片最大 2Mib
	imgToBase64SizeRate = 1.34    // 图片转base64 空间膨胀为 4/3

	robotApi = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
)

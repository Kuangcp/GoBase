package core

const (
	Open  = 1 // 开启配置
	Close = 0 // 关闭配置

	Direct = "DIRECT" // 直连
	Modify = "MOD"    // 抓包并修改
	Proxy  = "PROXY"  // 抓包代理

	// PacUrl https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
	PacUrl      = "/pac" // pac url
	pacFileType = "application/x-ns-proxy-autoconfig"
)

package core

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/kuangcp/logger"
	"net/http"
	"strings"
)

// 类似一个c++类，里面的属性初始化的时候是是可以赋值的
var upgrader = websocket.Upgrader{
	//此处给CheckOrigin默认一个返回true保证，否则会出现报错自动跳转
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func websocketHandler(w http.ResponseWriter, r *http.Request, proxyReq *http.Request) bool {
	if !strings.Contains(r.URL.Path, "js") {
		logger.Info(r.URL.Path, r.Header.Get("Upgrade"))
	}
	//此处检测到websocket请求
	if r.Header.Get("Upgrade") == "" {
		return false
	}

	/************进行协议升级************/
	upgrader.Subprotocols = []string{proxyReq.Header.Get("Sec-WebSocket-Protocol")}
	//upgrader.Upgrade内部会返回握手信息，我们做代理需要将dialer.Dial客户端收到的下层返回的握手信息返回给上层，源码Upgrade函数中将
	//Sec-WebSocket-Protocol这个头去掉了，所以我们给加上,保证在upgrader.Upgrade调用的上方加上
	//func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {}
	//upgrader.Subprotocols是属性，upgrader.Upgrade是方法，属性可以初始化，方法在函数定义前进行关联
	//具体实现在 https://github.com/gorilla/websocket/blob/master/server.go
	c_this, err := upgrader.Upgrade(w, proxyReq, nil)
	if err != nil {
		logger.Info("upgrade:", err)
		return true
	}
	defer c_this.Close()

	/******************启动websocket转发客户端*********/
	//此处使用req.URL.Path+req.URL.RawQuer代替 req.URL.String()是因为使用后者之后下面的u.String()使用url加密，导致url错误造成404
	//u := url.URL{Scheme: "wss", Host: proxyReq.Host, Path: proxyReq.Ul, RawQuery: req.URL.RawQuery}
	u := proxyReq.URL
	logger.Info("connecting to %s\n", u.String())
	//添加头部的时候不能照搬，需要去掉几个固有的，因为库里面已经给你加了，所以代理的时候把重复的去掉，详见源码219行
	//https://github.com/gorilla/websocket/blob/master/client.go
	headers := make(http.Header)
	for k, v := range proxyReq.Header {
		if k == "Upgrade" ||
			k == "Connection" ||
			k == "Sec-Websocket-Key" ||
			k == "Sec-Websocket-Version" ||
			k == "Sec-Websocket-Extensions" {
		} else {
			headers.Set(k, v[0])
			//fmt.Println("set ==>", k, v[0])
		}
	}
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	c_to_next, resp, err := dialer.Dial(u.String(), headers)
	if err != nil {
		logger.Info("dial:", err)
		logger.Info("StatusCode:", resp.StatusCode)
	}
	//fmt.Println(resp.Header)
	defer c_to_next.Close()
	/*****************接收返回并转回给浏览器*******************/

	go func() {
		logger.Info("run read from next proc")
		for {
			mt, message, err := c_to_next.ReadMessage()
			if err != nil {
				logger.Info("read from next:", err)
				break
			}
			//fmt.Println("read from next:", message)
			err = c_this.WriteMessage(mt, message)
			if err != nil {
				logger.Info("write to priv:", err)
				break
			}
		}
	}()

	/*****************接收浏览器信息并转发*******************/
	//此处不能再协程了，否则会defer c_to_next.Close()
	logger.Info("run read from priv proc")
	for {
		mt1, message1, err1 := c_this.ReadMessage()
		if err1 != nil {
			logger.Info("read from priv:", err1)
			break
		}
		//fmt.Println("read from priv:", message1)
		err1 = c_to_next.WriteMessage(mt1, message1)
		if err1 != nil {
			logger.Info("write to next:", err1)
			break
		}
	}
	/********************************************/
	return true
}

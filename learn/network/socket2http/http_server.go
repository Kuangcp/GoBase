package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func readRequest(conn *net.TCPConn) []byte {
	var maxBuf []byte
	bufSize := 1024
	var buf = make([]byte, bufSize)
	for {
		read, err := conn.Read(buf)
		//maxBuf = append(maxBuf, buf...) 如果这样写，会由于缓存没有清0而在最后一次读取时引入上一次内容
		maxBuf = append(maxBuf, buf[:read]...)
		if read != bufSize || err != nil {
			break
		}
	}
	return maxBuf
}

func buildAuthResponse() string {
	// Key:Value 的形式表示Header属性
	var header = []string{
		"HTTP/2.0 401 OK",
		"Server: Go Socket",
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/WWW-Authenticate
		"WWW-Authenticate: Basic realm=\"go Auth\", charset=\"UTF-8\"",
	}
	response := strings.Join(header, "\r\n")
	return response
}

func buildResponse() string {
	// Key:Value 的形式表示Header属性
	var header = []string{
		// HTTP版本(HTTP/1.0, HTTP/1.1) 和 状态(可以任意修改，因为只是规范)
		"HTTP/2.0 200 OK",
		"Server: Go Socket",
		// 用于避免tcp粘包问题，可以理解为定长解析器，
		// 客户端会严格按照此处声明长度进行读取,如果长于实际内容，浏览器或curl等客户端会不同形式报错，但是内容会正常展示，如果短了就会截断实际内容
		//"Content-Length: 4",
		"Connection: keep-alive",
	}
	response := strings.Join(header, "\r\n")

	// header 和 body 之间用一个换行 CRLF 分隔
	response += "\r\n\r\n"
	response = response + time.Now().Format("15:04:05.000")
	return response
}

func httpResponse(conn *net.TCPConn) {
	requestBuf := readRequest(conn)
	request := string(requestBuf)
	fmt.Println("----request----\n", request, "\n-------------")

	var response string
	if strings.Contains(request, "auth") {
		response = buildAuthResponse()
	} else {
		response = buildResponse()
	}

	n, err := conn.Write([]byte(response))
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
	conn.CloseWrite()
}

func main() {
	address := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"), // 把字符串IP地址转换为net.IP类型
		Port: 8000,
	}
	listener, err := net.ListenTCP("tcp4", &address) // 创建TCP4服务器端监听器
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err) // 错误直接退出
		}
		fmt.Println("remote address:", conn.RemoteAddr())
		go httpResponse(conn)
	}
}

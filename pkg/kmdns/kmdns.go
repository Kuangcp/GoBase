package kmdns

import (
	"context"
	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"time"
)

type (
	KmDNS struct {
		serviceName []string
		findTimeout time.Duration
	}
)

func New(findTimeout time.Duration, serviceName ...string) *KmDNS {
	return &KmDNS{serviceName: serviceName, findTimeout: findTimeout}
}

func (k *KmDNS) Server() {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}

	_, err = mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{
		LocalNames: k.serviceName,
	})
	if err != nil {
		panic(err)
	}
	select {}
}

// ClientRequest 客户端解析寻找服务端地址
func (k *KmDNS) ClientRequest() []string {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}

	server, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{})
	if err != nil {
		panic(err)
	}

	// 防止发出大量udp查询包 服务提供方不存在时，阻塞局域网
	timeout, cancelFunc := context.WithTimeout(context.TODO(), k.findTimeout)
	defer func() {
		cancelFunc()
		log.Println("cancel")
		//logger.Info("cancel")
	}()

	var res []string
	for _, name := range k.serviceName {
		answer, src, err := server.Query(timeout, name)
		if err != nil {
			log.Println(answer, src, err)
			continue
		}
		mdnsServer := src.String()
		res = append(res, mdnsServer[:len(mdnsServer)-4])
	}

	return res
}

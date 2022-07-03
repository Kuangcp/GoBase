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
		serviceName string
		timeout     time.Duration
	}
)

func New(serviceName string, timeout time.Duration) *KmDNS {
	return &KmDNS{serviceName: serviceName, timeout: timeout}
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
		LocalNames: []string{k.serviceName},
	})
	if err != nil {
		panic(err)
	}
	select {}
}

//FindMasterService
func (k *KmDNS) FindMasterService() string {
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
	timeout, cancelFunc := context.WithTimeout(context.TODO(), k.timeout)
	defer func() {
		cancelFunc()
		//logger.Info("cancel")
	}()

	answer, src, err := server.Query(timeout, k.serviceName)
	log.Println(answer, src, err)
	if err != nil {
		return ""
	}
	mdnsServer := src.String()
	return mdnsServer[:len(mdnsServer)-4]
}

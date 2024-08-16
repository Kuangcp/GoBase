package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net"
	"testing"
)

func TestIpNet(t *testing.T) {
	cidr, ipNet, err := net.ParseCIDR("192.168.1.0/24")
	if err != nil {
		logger.Error(err)
		return
	}

	fmt.Println(cidr, ipNet)
	start, end := AddressRange(ipNet)
	fmt.Println(start, end)
	cur := start
	for !cur.Equal(end) {
		fmt.Println(cur)
		cur = Inc(cur)
	}
}

func TestScanOneHostAll(t *testing.T) {
	hostStr = "192.168.9.174"
	//portStr = "7000-9000"
	allPort = true
	main()
}

func TestAllHost80Port(t *testing.T) {
	hostStr = "192.168.9.0/24"
	con = 50
	main()
}

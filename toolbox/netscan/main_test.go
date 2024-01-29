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

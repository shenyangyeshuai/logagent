package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
)

var (
	localIPs []string
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logs.Error(err)
		return
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPs = append(localIPs, ipnet.IP.String())
			}
		}
	}

	fmt.Println(localIPs)
}

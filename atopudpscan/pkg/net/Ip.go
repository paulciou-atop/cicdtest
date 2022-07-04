package net

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

//Get local IP list
func GetLocalIP() ([]string, error) {
	ips := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

//check ip format
func CheckIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("IP Address: %s - Invalid", ip)
	} else {
		return nil
	}
}

//check Mac format
func CheckMacAddress(mac string) error {
	_, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("MAC Address: %s - Invalid", mac)
	} else {
		return nil
	}
}

//check device exist
func CheckDeviceExisted(ip string) bool {
	r := false
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		return false
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		r = true
	}
	p.OnIdle = func() {
		log.Println("check finish")
	}
	err = p.Run()
	if err != nil {
		log.Println(err)
		return false
	}
	return r
}

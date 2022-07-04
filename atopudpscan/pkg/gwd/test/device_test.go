package gwd_test

import (
	"fmt"
	"log"
	"net"
	FirmWare "nms/atopudpscan/pkg/AtopFirmWare"
	"nms/atopudpscan/pkg/gwd"
	"os"
	"os/signal"
	"strings"
	"testing"
	"time"
)

func TestFwUpgrading(t *testing.T) {
	f, err := os.Open("EH7520-4G-8PoE-4SFP_K544_A544.dld")
	if err != nil {

		t.Error(err)
		return
	}
	defer f.Close()
	device := gwd.NewDevice("192.168.4.30").FirmWare()
	c := make(chan os.Signal, 1)
	go func() {
		for {
			time.Sleep(1000)
			s, err := device.GetProcessStatus()
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Print(s)
			if s == FirmWare.Complete || s == FirmWare.Error {
				break
			}

		}
	}()
	err = device.Upgrading(f)
	if err != nil {

		t.Error(err)
	}

	signal.Notify(c, os.Interrupt)
	<-c
}

func TestLogin(t *testing.T) {
	device := gwd.NewDevice("192.168.4.29")
	r, err := device.Login("default")
	if err != nil {

		t.Error(err)
	}
	if r {
		t.Log("pass")

	} else {
		t.Error("faield")

	}

}

func TestQuerySnmp(t *testing.T) {
	device := gwd.NewDevice("192.168.4.30")

	s, err := device.QuerySnmp("default")
	if err != nil {

		t.Error(err)
		return
	}
	fmt.Print(s.Name)
	fmt.Print(s.Location)
	fmt.Print(s.Contact)
}

func TestSettingSnmp(t *testing.T) {
	device := gwd.NewDevice("192.168.4.30")
	snmp := &gwd.Snmp{Contact: "test", Name: "52t789", Location: "taivh"}

	s, err := device.SettingSnmp("default", snmp)
	if err != nil {

		t.Error(err)
		return
	}
	fmt.Print(s.Name)
	fmt.Print(s.Location)
	fmt.Print(s.Contact)
}

func TestLocalIp(t *testing.T) {

	ips, err := GetLocalIP()
	if err != nil {
		t.Error(err)
	}
	stringByte := strings.Join(ips, "\x2c")
	fmt.Println([]byte(stringByte))
	fmt.Println(string([]byte(stringByte)))

}
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

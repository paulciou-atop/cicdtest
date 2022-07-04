package gwd_test

import (
	"nms/atopudpscan/pkg/gwd"
	"sync"
	"testing"
)

func TestScanner(t *testing.T) {

	g := gwd.NewAtopGwd("0.0.0.0:0")
	err := g.Scanner()
	if err != nil {
		t.Error(err)
	}

}
func TestMutipleScanner(t *testing.T) {
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		go func() {
			wg.Add(1)
			g := gwd.NewAtopGwd("192.168.4.21:0")
			err := g.Scanner()
			if err != nil {
				t.Error(err)
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func TestBeep(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:  "192.168.4.30",
		MACAddress: "00:60:e9:18:3c:3c",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.Beep(net)
	if err != nil {

		t.Error(err)
	}
}

func TestReboot(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:  "192.168.4.25",
		MACAddress: "00:60:e9:11:22:33",
		Username:   "admin",
		Password:   "default",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.Reboot(net)
	if err != nil {

		t.Error(err)
	}
}

func TestResetToDefault(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:  "192.168.4.30",
		MACAddress: "00:60:e9:11:22:33",
		Username:   "admin",
		Password:   "default",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.ResetToDefault(net)
	if err != nil {

		t.Error(err)
	}
}

func TestSettingConfig(t *testing.T) {
	net := gwd.NetworkConfig{
		IPAddress:    "192.168.4.30",
		MACAddress:   "00:60:e9:18:3c:3c",
		NewIPAddress: "192.168.4.30",
		Netmask:      "255.255.255.0",
		Gateway:      "192.168.4.254",
		Hostname:     "atoptest",
		Username:     "admin",
		Password:     "default",
	}

	g := gwd.NewAtopGwd("192.168.4.21:0")
	err := g.SettingConfig(net)
	if err != nil {

		t.Error(err)
	}
}

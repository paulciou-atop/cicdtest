package Simulate_test

import (
	"log"
	"nms/atopudpscan/pkg/Simulatedevice"
	"os"
	"os/signal"
	"testing"
)

const number = 500
const ip = "192.168.4.22"

func TestServer(t *testing.T) {
	s := Simulatedevice.NewSimulateGwdServer(ip)
	for i := 1; i <= number; i++ {
		d, err := Simulatedevice.NewAtopDevice(ip, "device", uint(i))
		if err != nil {
			log.Fatal(err)
		}
		s.RegisterHandler(d)
	}
	err := s.Run()
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func TestParma(t *testing.T) {
	m, err := Simulatedevice.GetTestParam("test", 1)
	if err != nil {
		t.Error(err)
	}

	log.Print(m)
	//{testdevice1 00-60-E9-2C-00-01 192.168.1.1 255.255.255.0 192.168.1.254 test1   false}
}

package Simulatedevice

import (
	"bytes"
	"errors"
	"log"
	"net"
	"nms/atopudpscan/pkg/gwd"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Create  Atop Device
//
//Number: sequence on ip, Hostname
//
//Exampe:Hostname=device,number=1
//
//Model=testdevice
//
//MACAddress: "00-60-E9-2C-01:"+nmuber = 00-60-E9-2C-01-01
//
//IPAddress:"192.168.9." + number = 192.168.9.1
//
//Netmask: "255.255.255.0"
//Gateway: "192.168.9.254"
//
//Hostname:"device" + number =device1
func NewAtopDevice(outip, Hostname string, number uint) (*AtopGwdClient, error) {
	Model, _ := GetTestParam(Hostname, number)
	mac := strings.Split(Model.MACAddress, "-")
	macb := make([]byte, 6)
	for i := 0; i < len(mac); i++ {
		v, _ := strconv.ParseUint(mac[i], 16, 8)
		macb[i] = byte(v)
	}
	d := &AtopGwdClient{ModelInfo: Model, mac: macb, l: new(sync.Mutex), startup: true, outip: outip}
	err := checkLenHost(d.ModelInfo.Hostname)
	return d, err
}

type AtopGwdClient struct {
	ModelInfo gwd.ModelInfo
	mac       []byte
	l         *sync.Mutex
	startup   bool
	outip     string
}

func (a *AtopGwdClient) SendDevice() {
	a.l.Lock()
	a.broadcast(a.ModelInfoPacket())
	a.l.Unlock()
}

func checkLenHost(name string) error {
	if len(name) > 15 {
		return errors.New("Name len than 15")
	}
	return nil
}
func (a *AtopGwdClient) SettingDevice(msg []byte) {
	if a.compareMac(msg) {
		a.l.Lock()
		newip := make([]string, 4)
		for i := 0; i < 4; i++ {
			newip[i] = strconv.Itoa(int(msg[16+i]))
		}
		newmask := make([]string, 4)
		for i := 0; i < 4; i++ {
			newmask[i] = strconv.Itoa(int(msg[236+i]))
		}

		gate := make([]string, 4)
		for i := 0; i < 4; i++ {
			gate[i] = strconv.Itoa(int(msg[24+i]))
		}

		hostname := make([]string, 15)
		for i := 0; i < len(hostname); i++ {
			hostname[i] = string(msg[90+i])
		}
		a.ModelInfo.IPAddress = strings.Join(newip[:], ".")
		a.ModelInfo.Netmask = strings.Join(newmask[:], ".")
		a.ModelInfo.Gateway = strings.Join(gate[:], ".")
		n := strings.Join(hostname[:], "")
		err := checkLenHost(n)
		if err == nil {
			a.ModelInfo.Hostname = n
		}
		a.startup = false
		a.l.Unlock()
		a.Reboot()
	}

}

func (a *AtopGwdClient) Setstatus(s bool) {
	a.l.Lock()
	a.startup = s
	a.l.Unlock()
}

func (a *AtopGwdClient) GetStatus() bool {
	a.l.Lock()
	v := a.startup
	a.l.Unlock()
	return v
}

func (a *AtopGwdClient) Reboot() {
	a.Setstatus(false)
	c := make(chan bool)
	go func() {
		c <- true
		time.Sleep(time.Second * 10)
		a.Setstatus(true)
	}()
	<-c
}

func (a *AtopGwdClient) broadcast(msg []byte) {
	addr := net.JoinHostPort("255.255.255.255", strconv.Itoa(gwd.Port))
	broadcastAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Print(err)
	}
	local, err := net.ResolveUDPAddr("udp", net.JoinHostPort(a.outip, strconv.Itoa(0)))
	if err != nil {
		log.Print(err)
	}
	conn, err := net.DialUDP("udp", local, broadcastAddr)
	if err != nil {
		log.Print(err)
	}
	_, err = conn.Write(msg)
	if err != nil {
		log.Print(err)
	}
}

func (a *AtopGwdClient) Receive(b []byte, s *net.UDPAddr) {
	r := SelectPacket(b)
	switch r {
	case invite:
		if a.GetStatus() {
			a.SendDevice()
		}
	case config:
		a.SettingDevice(b)
	case none:
	}
}

func (a *AtopGwdClient) compareMac(m []byte) bool {
	r := bytes.Compare(a.mac, m[28:34])
	if r == 0 {
		return true
	}
	return false
}
func (a *AtopGwdClient) ModelInfoPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 0x01
	packet[4] = 0x92
	packet[5] = 0xDA
	modinfo := []byte(a.ModelInfo.Model)
	for i := 0; i < len(modinfo); i++ {
		packet[44+i] = modinfo[i]
	}
	ip := strings.Split(a.ModelInfo.IPAddress, ".")
	for i := 0; i < len(ip); i++ {
		v, _ := strconv.Atoi(ip[i])
		packet[12+i] = byte(v)
	}
	mac := strings.Split(a.ModelInfo.MACAddress, "-")
	for i := 0; i < len(mac); i++ {
		v, _ := strconv.ParseUint(mac[i], 16, 8)
		packet[28+i] = byte(v)
	}
	Netmask := strings.Split(a.ModelInfo.Netmask, ".")
	for i := 0; i < len(Netmask); i++ {
		v, _ := strconv.Atoi(Netmask[i])
		packet[236+i] = byte(v)
	}

	Gateway := strings.Split(a.ModelInfo.Gateway, ".")
	for i := 0; i < len(Gateway); i++ {
		v, _ := strconv.Atoi(Gateway[i])
		packet[24+i] = byte(v)
	}
	Hostname := []byte(a.ModelInfo.Hostname)
	for i := 0; i < len(Hostname); i++ {
		packet[90+i] = Hostname[i]
	}

	return packet
}

func SelectPacket(packet []byte) int {
	if packet[0] == 0x02 && packet[1] == 0x01 && packet[2] == 0x06 && packet[4] == 0x92 && packet[5] == 0xDA {
		return invite
	}
	if packet[0] == 0 && packet[1] == 1 && packet[2] == 6 && packet[4] == 0x92 && packet[5] == 0xDA {
		mac := make([]string, 6)
		for i := 0; i < 6; i++ {
			mac[i] = string(packet[28])
		}
		return config

	}
	return none
}

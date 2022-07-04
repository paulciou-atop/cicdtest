package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"nms/atopudpscan/pkg/gwd"
	"os"
	"strconv"
	"strings"
	"sync"
)

const fileName = "AtopGwd.Log"

func NewAtopUdpServer(ip string) *AtopGwdServer {
	return &AtopGwdServer{ip: ip, deviceList: make([]gwd.ModelInfo, 0),
		m: new(sync.Mutex), listening: false}
}

type AtopGwdServer struct {
	ip         string
	deviceList []gwd.ModelInfo
	conn       *net.UDPConn
	m          *sync.Mutex
	listening  bool
}

//udp server run
func (r *AtopGwdServer) Run() error {
	address := strings.Join([]string{r.ip, strconv.Itoa(gwd.Port)}, ":")
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {

		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {

		return err
	}
	log.Printf("start listen...%s\n", address)

	if err != nil {
		log.Println("read from connect failed, err:" + err.Error())
		return err
	}
	r.conn = conn
	c := make(chan bool, 1)
	go r.udpProcess(conn, c)
	<-c
	return nil
}

func (r *AtopGwdServer) Scan(outIp string) error {
	r.m.Lock()
	r.deviceList = nil
	g := gwd.NewAtopGwd(outIp)
	err := g.Scanner()
	r.m.Unlock()
	if err != nil {
		return err
	}
	return nil
}

//get json byte date from udp storge
func (r *AtopGwdServer) GetReceiveData() ([]byte, error) {
	r.m.Lock()
	v := r.deviceList
	r.m.Unlock()
	if len(v) == 0 {
		return nil, errors.New("len:0")
	}
	jsondata, err := json.Marshal(v)

	return jsondata, err
}

func (r *AtopGwdServer) ListeningStatus() bool {
	r.m.Lock()
	b := r.listening
	r.m.Unlock()
	return b
}

func (r *AtopGwdServer) SaveToDatabase(b []byte) error {
	return writeFile(fileName, b)
}

func (r *AtopGwdServer) GetDataFromDatabase() ([]byte, error) {
	return ioutil.ReadFile(fileName)

}
func (r *AtopGwdServer) SearchDeviceMAC(mac string) ([]byte, error) {
	r.m.Lock()
	v := r.deviceList
	r.m.Unlock()
	if len(v) == 0 {
		return nil, errors.New("len:0")
	}
	var d *gwd.ModelInfo
	for _, v := range v {
		if v.MACAddress == mac {
			dv := v
			d = &dv
			break
		}
	}
	if d == nil {
		return nil, fmt.Errorf("MACAddress:%v is not exist", mac)
	}
	jsondata, err := json.Marshal(d)
	return jsondata, err
}

func (r *AtopGwdServer) CleanDatabase() error {
	if err := os.Truncate(fileName, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
		return nil
	} else {
		return err
	}
}

func (r *AtopGwdServer) udpProcess(conn *net.UDPConn, c chan bool) {
	r.m.Lock()
	r.listening = true
	c <- true
	r.m.Unlock()
	for {
		data := make([]byte, 300)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			//log.Println("failed read udp msg, error: " + err.Error())
			r.m.Lock()
			r.listening = false
			r.m.Unlock()
			return
		}
		//str := string(data[:n])
		info, err := parsingModelInfo(data[:n])
		r.m.Lock()
		if err == nil {
			r.deviceList = append(r.deviceList, *info)
		}
		r.m.Unlock()
		//jsondata, _ := json.Marshal(info)
		//fmt.Println("receive from client, data:" + string(jsondata))
	}

}
func (r *AtopGwdServer) Stop() error {
	log.Printf("server stop")

	if r.conn != nil {
		err := r.conn.Close()
		r.conn = nil
		return err
	}
	return errors.New("No scan running")
}

//parsing data form udp
func parsingModelInfo(msg []byte) (*gwd.ModelInfo, error) {
	if len(msg) != 300 {
		return nil, errors.New(fmt.Sprint("len is error:", len(msg)))
	}
	if msg[0] == 1 && msg[4] == 0x92 && msg[5] == 0xDA {
		model := &gwd.ModelInfo{}
		model.Model = toUtf8(msg[44:60])
		model.MACAddress = byteToHexString(msg[28:34], "-")
		model.IPAddress = byteToString(msg[12:16], ".")
		model.Netmask = byteToString(msg[236:240], ".")
		model.Gateway = byteToString(msg[24:28], ".")
		model.Hostname = toUtf8(msg[90:106])
		model.Kernel = fmt.Sprintf("%d.%d", msg[109], msg[108])
		model.Ap = toUtf8(msg[110:235])
		if fmt.Sprintf("%d", msg[106]) == "0" {
			model.IsDHCP = false
		} else {
			model.IsDHCP = true
		}
		//jsondata, _ := json.Marshal(model)
		return model, nil
	} else {
		return nil, errors.New("error")
	}
}

func writeFile(filename string, value []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.Write(value)
	f.Write([]byte("\n"))
	defer f.Close()
	return nil
}

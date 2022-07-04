package Simulatedevice

import (
	"errors"
	"fmt"
	"log"
	"net"
	"nms/atopudpscan/pkg/gwd"
	"strconv"
	"strings"
	"sync"
)

const none = 0
const invite = 1
const config = 2

type GwdDevice interface {
	Receive([]byte, *net.UDPAddr)
}

//Create New SimulaterGwd
//
//localhost: listenting IP
func NewSimulateGwdServer(localhost string) *SimulateGwdServer {
	return &SimulateGwdServer{ip: localhost, deviceList: make([]gwd.ModelInfo, 0), m: new(sync.Mutex)}
}

type SimulateGwdServer struct {
	ip         string
	conn       *net.UDPConn
	m          *sync.Mutex
	hanlder    []GwdDevice
	deviceList []gwd.ModelInfo
}

const SimulatePort = 55954

func (r *SimulateGwdServer) Run() error {
	address := strings.Join([]string{r.ip, strconv.Itoa(SimulatePort)}, ":")
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

func (r *SimulateGwdServer) udpProcess(conn *net.UDPConn, c chan bool) {
	r.m.Lock()
	c <- true
	r.m.Unlock()
	for {
		data := make([]byte, 300)
		n, soruce, err := conn.ReadFromUDP(data)
		if err != nil {
			//log.Println("failed read udp msg, error: " + err.Error())
			log.Print(err)
			return
		}
		//str := string(data[:n])
		r.m.Lock()
		//	info, err := parsingModelInfo(data[:n])
		//if err == nil {
		//		r.deviceList = append(r.deviceList, *info)
		//		log.Print(info)
		//	} else {
		for _, v := range r.hanlder {
			v.Receive(data[:n], soruce)
		}
		//	}
		r.m.Unlock()
		//jsondata, _ := json.Marshal(info)
		//fmt.Println("receive from client, data:" + string(jsondata))
	}

}

//Register behavior of device Receive
func (r *SimulateGwdServer) RegisterHandler(g GwdDevice) {
	r.m.Lock()
	r.hanlder = append(r.hanlder, g)
	r.m.Unlock()
}

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

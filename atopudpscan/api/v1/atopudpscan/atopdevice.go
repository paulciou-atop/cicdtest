package atopudpscan

import (
	context "context"
	"errors"
	"fmt"
	"io"
	"log"
	"nms/atopudpscan/internal/pkg/file"
	"nms/atopudpscan/pkg/gwd"
	"nms/atopudpscan/pkg/net"
	"os"
	"strconv"
	"strings"
	sync "sync"
	"time"

	api "nms/api/v1/atopudpscan"
	FirmWare "nms/atopudpscan/pkg/AtopFirmWare"
)

const waitsettingtimeout = 120

func NewDeviceController() *DeviceController {

	d := &DeviceController{FwdeviceList: make([]*gwd.Device, 0), fwLock: new(sync.Mutex)}
	AtopudpscanShare().RegisterDeviceController(d)
	return d
}

type DeviceController struct {
	fwLock       *sync.Mutex
	FwdeviceList []*gwd.Device
	api.UnimplementedAtopDeviceServer
}

//if device is upgrade, return device info
func (d *DeviceController) searchFwDevice(ip string) (int, *gwd.Device) {
	d.fwLock.Lock()
	f := d.FwdeviceList
	d.fwLock.Unlock()
	for i, v := range f {
		if v.Ip() == ip {

			return i, v
		}

	}

	return 0, nil
}

//set device config
func (d *DeviceController) SettingConfig(ctx context.Context, in *api.GwdConfig) (*api.Response, error) {

	_, a := d.searchFwDevice(in.GetIPAddress())
	if a != nil {
		s, _ := a.FirmWare().GetProcessStatus()
		if s != FirmWare.Complete && s != FirmWare.Error {
			return &api.Response{Result: false, Message: "device is Upgrading"}, fmt.Errorf("device is Upgrading")
		}
	}

	net := gwd.NetworkConfig{
		IPAddress:    in.GetIPAddress(),
		MACAddress:   in.GetMACAddress(),
		NewIPAddress: in.GetNewIPAddress(),
		Netmask:      in.GetNetmask(),
		Gateway:      in.GetGateway(),
		Hostname:     in.GetHostname(),
		Username:     in.GetUsername(),
		Password:     in.GetPassword(),
	}
	address := strings.Join([]string{in.GetServerIp(), strconv.Itoa(0)}, ":")
	g := gwd.NewAtopGwd(address)
	err := g.SettingConfig(net)
	if err != nil {
		return &api.Response{Result: false, Message: err.Error()}, err
	}
	return &api.Response{Result: true}, err

}

/*
func (d *DeviceController) ResetToDefault(ctx context.Context, in *GwdConfig) (*Response, error) {
	_, a := d.searchFwDevice(in.GetIPAddress())
	if a != nil {
		return &Response{Result: false, Message: "device is Upgrading"}, fmt.Errorf("device is Upgrading")
	}

	net := gwd.NetworkConfig{
		IPAddress:  in.GetIPAddress(),
		MACAddress: in.GetMACAddress(),
		Username:   in.GetUsername(),
		Password:   in.GetPassword(),
	}
	address := strings.Join([]string{in.GetServerIp(), strconv.Itoa(0)}, ":")
	g := gwd.NewAtopGwd(address)
	err := g.ResetToDefault(net)
	if err != nil {
		return &Response{Result: false, Message: err.Error()}, err
	}
	return &Response{Result: true}, nil

}*/

//reboot device
func (d *DeviceController) Reboot(ctx context.Context, in *api.GwdConfig) (*api.Response, error) {
	_, a := d.searchFwDevice(in.GetIPAddress())
	if a != nil {
		s, _ := a.FirmWare().GetProcessStatus()
		if s != FirmWare.Complete && s != FirmWare.Error {
			return &api.Response{Result: false, Message: "device is Upgrading"}, fmt.Errorf("device is Upgrading")
		}
	}

	net := gwd.NetworkConfig{
		IPAddress:  in.GetIPAddress(),
		MACAddress: in.GetMACAddress(),
		Username:   in.GetUsername(),
		Password:   in.GetPassword(),
	}
	address := strings.Join([]string{in.GetServerIp(), strconv.Itoa(0)}, ":")
	g := gwd.NewAtopGwd(address)
	err := g.Reboot(net)
	if err != nil {
		return &api.Response{Result: false, Message: err.Error()}, err
	}
	return &api.Response{Result: true}, nil

}

//device fw upgrade after complete file deleted
func (d *DeviceController) FwUpgrading(ctx context.Context, in *api.FwInfo) (*api.Response, error) {
	var atop *gwd.Device

	deviceip := in.GetDeviceIp()

	err := net.CheckIPAddress(deviceip)
	if err != nil {
		log.Print(err)
		return &api.Response{Result: false, Message: err.Error()}, err
	}

	f1, err := os.Open(in.FileName)
	if err != nil {
		log.Print(err)
		return &api.Response{Result: false, Message: err.Error()}, err
	}
	//upload fw
	_, a := d.searchFwDevice(deviceip)
	if a != nil {
		s, _ := a.FirmWare().GetProcessStatus()
		if s != FirmWare.Complete && s != FirmWare.Error {
			return &api.Response{Result: false, Message: "device is Upgrading"}, fmt.Errorf("device is Upgrading")
		}
		atop = a
	} else {
		atop = gwd.NewDevice(deviceip)
		d.add(atop)
	}

	fw := atop.FirmWare()
	err = fw.Upgrading(f1)
	if err != nil {
		log.Print(err)
		return &api.Response{Result: false, Message: err.Error()}, err
	}

	c := make(chan bool, 1)
	go func() {
		c <- true
		//remove file after complete or error
		//defer removeAll(in.FileName)
		defer f1.Close()
		for {
			s, err := fw.GetProcessStatus()
			if err != nil {
				return
			}
			if s == FirmWare.Complete || s == FirmWare.Error {
				return
			}
		}

	}()
	<-c
	return &api.Response{Result: true}, nil
}

//upload file and return new filename
func (d *DeviceController) Upload(stream api.AtopDevice_UploadServer) error {
	fwname := time.Now().Format("20060102150405") + ".dld"
	f := file.NewFile(fwname)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			if err := f.Save(); err != nil {
				log.Print(err)
				stream.SendAndClose(&api.UploadStatus{Code: api.UploadStatusCode_Failed, Message: err.Error()})
				return err
			}

			break
		}
		if err != nil {
			log.Printf("cannot receive %v", err)
			stream.SendAndClose(&api.UploadStatus{Code: api.UploadStatusCode_Failed, Message: err.Error()})
			return err
		}

		if err := f.Write(resp.GetContent()); err != nil {
			stream.SendAndClose(&api.UploadStatus{Code: api.UploadStatusCode_Failed, Message: err.Error()})
			log.Print(err)
			return err
		}

	}
	stream.SendAndClose(&api.UploadStatus{Code: api.UploadStatusCode_Ok, FileName: fwname})
	return nil
}

func (d *DeviceController) add(atop *gwd.Device) {
	d.fwLock.Lock()
	d.FwdeviceList = append(d.FwdeviceList, atop)
	d.fwLock.Unlock()
}

//get upgrading status of device
func (d *DeviceController) GetProcessStatus(ctx context.Context, in *api.FwRequest) (*api.FwMessage, error) {
	_, a := d.searchFwDevice(in.GetIPAddress())
	if a == nil {
		return &api.FwMessage{Message: "device not exist"}, errors.New("device not exist")
	}
	v, err := a.FirmWare().GetProcessStatus()
	if err == nil {
		switch v {
		case FirmWare.None:
			return &api.FwMessage{Status: api.FwStatus_none}, err
		case FirmWare.Complete:
			return &api.FwMessage{Status: api.FwStatus_complete}, err
		case FirmWare.Upgrading:
			return &api.FwMessage{Status: api.FwStatus_upgrading}, err
		default:
			return &api.FwMessage{Status: api.FwStatus_process, Message: v}, err
		}
	} else {
		return &api.FwMessage{Status: api.FwStatus_error, Message: err.Error()}, err
	}
}

//get server ips
func (d *DeviceController) GetServerIp(ctx context.Context, in *api.EmptyParams) (*api.ServerIp, error) {
	ips, err := net.GetLocalIP()
	if err != nil {
		return &api.ServerIp{}, err
	}
	return &api.ServerIp{Ip: ips}, err
}

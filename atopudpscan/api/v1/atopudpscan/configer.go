package atopudpscan

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	api "nms/api/v1/atopudpscan"
	"nms/api/v1/configer"
	"nms/atopudpscan/internal/pkg/file"
	FirmWare "nms/atopudpscan/pkg/AtopFirmWare"
	"nms/atopudpscan/pkg/gwd"
	"nms/atopudpscan/pkg/net"

	"nms/api/v1/devconfig"

	"github.com/fatih/structs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

const gwdproto = "gwd"

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	configer.UnimplementedConfigerServer
}

//Valid date
func (c *Config) Validate(stream configer.Configer_ValidateServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		//:var failConfigs = []*devconfig.FailConfig{}
		if req == nil {
			return fmt.Errorf("null request")
		}
		if req.Configs == nil {
			return fmt.Errorf("missing config")
		}
		if req.Session == nil {
			return fmt.Errorf("missing session")
		}

		configs := req.Configs
		session := req.Session
		configResults := checkGwdUnSupport(configs) //check unsupport list

		if len(configResults) > 0 {
			session = sessionFail(session, fmt.Sprintf("some config not support"))
		} else {
			session = sessionSuccess(session)
		}
		stream.Send(&configer.ValidateResponse{
			Session:       session,
			ConfigResults: configResults,
		})

		return nil
	}
}

//check unSupport return unSupport array
func checkGwdUnSupport(configs []*devconfig.ConfigOptions) []*devconfig.ConfigResult {
	var configResults = []*devconfig.ConfigResult{}
	for _, config := range configs {
		if config.Protocol != gwdproto { //only gwd
			configResults = append(configResults, &devconfig.ConfigResult{
				Protocol:   config.Protocol,
				Kind:       config.Kind,
				Hash:       config.Hash,
				FailFields: []string{"only support gwd"},
			})
			return configResults
		}
		unsupported := checkSupportConfig(config.Payload.AsMap())

		if len(unsupported) > 0 {
			configResults = append(configResults, &devconfig.ConfigResult{
				Protocol:   config.Protocol,
				Kind:       config.Kind,
				Hash:       config.Hash,
				FailFields: unsupported,
			})
		}
	}
	return configResults
}

//FileTransfer
func (c *Config) FileTransfer(stream configer.Configer_FileTransferServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		//get session
		session := resp.GetSession()
		if session == nil {
			return errors.New("Session is nil")
		}
		if session.Id == "" {
			return errors.New("id is empty")
		}
		f := resp.GetFileParam()
		if f == nil {
			return errors.New("FileParam is nil")
		}
		if f.Device.DevicePath == "" { //check ip
			return errors.New("DevicePath empty")
		}
		err = net.CheckIPAddress(f.Device.DevicePath) //check ip format
		if err != nil {

			return err
		}
		if f.FileLocation == "" {
			return errors.New("FileLocation is empty")
		}
		fwname := time.Now().Format("20060102150405") + ".dld" //creat time as filename
		err = DownloadFile(fwname, f.FileLocation)
		if err != nil {

			return err
		}
		f1, err := os.Open(fwname)
		if err != nil {
			return err
		}
		var atop *gwd.Device
		d := AtopudpscanShare().GetDeviceController()
		_, a := d.searchFwDevice(f.Device.DevicePath) //check device stats if it is upgrading
		if a != nil {
			s, _ := a.FirmWare().GetProcessStatus()
			if s != FirmWare.Complete && s != FirmWare.Error {
				return fmt.Errorf("%v is %v", f.Device.DevicePath, s)
			}
			atop = a
		} else {
			atop = gwd.NewDevice(f.Device.DevicePath)
			d.add(atop)
		}
		fw := atop.FirmWare()
		err = fw.Upgrading(f1)
		if err != nil {
			return err
		}
		session.State = "running"
		conf := &configer.ConfigerFileTransferResponse{Session: session, ByteTransfered: 0, FilePaht: f.FileLocation}
		err = stream.Send(conf)
		if err != nil {
			return err
		}
		defer file.Remove(fwname)
		defer f1.Close()
		//remove file after complete or error
		//defer removeAll(in.FileName)
		for {
			s, err := fw.GetProcessStatus()
			if err != nil {
				return err
			}
			if s == FirmWare.Complete {
				session.State = "success"
				conf = &configer.ConfigerFileTransferResponse{Session: session, ByteTransfered: 0, FilePaht: f.FileLocation}
				err = stream.Send(conf)
				if err != nil {
					return err
				}
				return nil
			}
			if s == FirmWare.Error {
				session.State = "fail "
				conf = &configer.ConfigerFileTransferResponse{Session: session, ByteTransfered: 0, FilePaht: f.FileLocation}
				err = stream.Send(conf)
				if err != nil {
					return err
				}
				return nil
			}
		}

	}

}

//configure
func (c *Config) Config(stream configer.Configer_ConfigServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if req == nil {
			return fmt.Errorf("null request")
		}
		if req.Configs == nil {
			return fmt.Errorf("missing config")
		}
		if req.Session == nil {
			return fmt.Errorf("missing session")
		}
		if req.Device == nil {
			return fmt.Errorf("missing device")
		}

		configs := req.Configs

		configResults := checkGwdUnSupport(configs)

		if len(configResults) > 0 {
			session := sessionFail(req.Session, fmt.Sprintf("some config not support"))
			stream.Send(&configer.ConfigerResponse{
				Session:       session,
				Device:        req.Device,
				ConfigResults: configResults,
			})
			return nil
		}

		//=======get original device value and search device==============//

		device := req.Device
		d, err := ScanandGetDeviceByMac("0.0.0.0", device.DeviceId)
		if err != nil {
			session := req.Session
			session = sessionFail(session, fmt.Sprintf("can not sreach device %s", device.DeviceId))
			err = stream.Send(&configer.ConfigerResponse{
				Session: session,
				Device:  req.Device,
			})
			if err != nil {
				return err
			}
			continue
		}

		//=======end get original device value and search device==============//

		//=========check device is upgrade fw=============//
		_, a := AtopudpscanShare().GetDeviceController().searchFwDevice(d.IPAddress)
		if a != nil {
			s, _ := a.FirmWare().GetProcessStatus()
			if s != FirmWare.Complete && s != FirmWare.Error {
				return fmt.Errorf("device is Upgrading")
			}
		}
		//=========endcheck device is upgrade fw=============//

		//===========check Protocol====================//
		var gwdconf *api.GwdConfig
		var conf *devconfig.ConfigOptions
		for _, config := range req.Configs {
			if config.Protocol == gwdproto {
				conf = config
				gwdconf, err = parsingConfig(config)
				if err != nil {
					session := req.Session
					session = sessionFail(session, err.Error())
					err = stream.Send(&configer.ConfigerResponse{
						Session: session,
						Device:  req.Device,
					})
					if err != nil {
						return err
					}
					return nil
				}
				break
			}
		}

		//===========end check Protocol====================//
		//===find device and Replenish Setting config if value is null===//
		gwdconf = ReplenishSettingConfig(d, gwdconf)
		address := strings.Join([]string{gwdconf.GetServerIp(), strconv.Itoa(0)}, ":")
		g := gwd.NewAtopGwd(address)
		net := gwd.NetworkConfig{
			IPAddress:    gwdconf.GetIPAddress(),
			MACAddress:   device.DeviceId,
			NewIPAddress: gwdconf.GetNewIPAddress(),
			Netmask:      gwdconf.GetNetmask(),
			Gateway:      gwdconf.GetGateway(),
			Hostname:     gwdconf.GetHostname(),
			Username:     gwdconf.GetUsername(),
			Password:     gwdconf.GetPassword(),
		}
		err = g.SettingConfig(net) //device set
		if err != nil {
			session := req.Session
			session = sessionFail(session, err.Error())
			err = stream.Send(&configer.ConfigerResponse{
				Session: session,
				Device:  req.Device,
			})
			if err != nil {
				return err
			}
		}
		session := req.Session
		session.State = "running"
		err = stream.Send(&configer.ConfigerResponse{
			Session: session,
			Device:  req.Device,
		})
		if err != nil {
			return err
		}
		//===end find device and Replenish Setting config if value is null===//
		time.Sleep(time.Second * 1)
		//===check deivce info ===//
		b, fail := checkDeviceIsChanged(gwdconf.GetServerIp(), net, waitsettingtimeout)
		if b {

			session := sessionSuccess(req.Session)
			err = stream.Send(&configer.ConfigerResponse{
				Session: session,
				Device:  req.Device,
			})
			if err != nil {
				return err
			}
		} else {
			var configResults = []*devconfig.ConfigResult{}
			res := devconfig.ConfigResult{Protocol: conf.Protocol, Kind: conf.Kind, Hash: conf.Hash, FailFields: fail}
			configResults = append(configResults, &res)
			var session *devconfig.SessionState
			if fail == nil {
				session = sessionFail(req.Session, "compare error")
			} else {
				session = sessionFail(req.Session, "not searh device")
			}
			err = stream.Send(&configer.ConfigerResponse{
				Session:       session,
				Device:        req.Device,
				ConfigResults: configResults,
			})
			if err != nil {
				return err
			}
		}
		return nil

		//===end check deivce info ===//
	}

}

//Get config
func (c *Config) GetConfig(ctx context.Context, req *configer.GetConfigRequest) (*configer.GetConfigResponse, error) {
	d, err := ScanandGetDeviceByMac("0.0.0.0", req.Device.DeviceId)
	if err != nil {
		return &configer.GetConfigResponse{}, err
	}
	b, err := json.Marshal(d)
	if err != nil {
		return &configer.GetConfigResponse{}, err
	}
	p := &structpb.Struct{}
	err = protojson.Unmarshal(b, p)
	if err != nil {
		return &configer.GetConfigResponse{}, err
	}
	return &configer.GetConfigResponse{
		Device:  req.Device,
		Configs: p,
	}, nil
}

//check Support config
func checkSupportConfig(payload map[string]interface{}) []string {
	var unsupportList = []string{}
	for field, _ := range payload {
		if !contains(GetGwdSupportConfig(), field) {
			unsupportList = append(unsupportList, field)
		}
	}
	return unsupportList
}

// contains check slice s contains e
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//get supportlist
func GetGwdSupportConfig() []string {
	m := make([]string, 0)
	g := api.GwdConfig{}
	f := structs.Map(&g)
	for k, _ := range f {
		m = append(m, k)
	}
	return m
}

func parsingConfig(c *devconfig.ConfigOptions) (*api.GwdConfig, error) {
	d := c.Payload
	b, err := protojson.Marshal(d)
	if err != nil {
		return nil, err
	}
	conf := &api.GwdConfig{}
	err = protojson.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	err = checkMustExistArgs(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

//check MustExistArgs
func checkMustExistArgs(c *api.GwdConfig) error {
	if c.ServerIp == "" {
		return fmt.Errorf("ServerIp must exist")
	}

	if c.Username == "" {
		return fmt.Errorf("Username must exist")
	}
	if c.Password == "" {
		return fmt.Errorf("Password must exist")
	}

	return nil
}

//scan device by mac
func ScanandGetDeviceByMac(outip, mac string) (*api.DeviceInfo, error) {
	r := AtopudpscanShare().GetGwd()
	address := strings.Join([]string{outip, strconv.Itoa(0)}, ":")
	err := r.u.Scan(address)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 1)

	v, err := r.u.SearchDeviceMAC(strings.ToUpper(mac))
	if err != nil {
		return nil, err
	}
	var data *api.DeviceInfo
	err = json.Unmarshal(v, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ReplenishSettingConfig(d *api.DeviceInfo, g *api.GwdConfig) *api.GwdConfig {
	c := g
	if c.IPAddress == "" {
		c.IPAddress = d.IPAddress
	}
	if c.Netmask == "" {
		c.Netmask = d.Netmask

	}
	if c.Gateway == "" {
		c.Gateway = d.Gateway
	}

	if c.NewIPAddress == "" {
		c.NewIPAddress = d.IPAddress
	}
	if c.Hostname == "" {
		c.Hostname = d.Hostname
	}

	return c
}

//check device is chang in timeout
func checkDeviceIsChanged(out string, g gwd.NetworkConfig, timeout int) (bool, []string) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	for {
		select {
		case <-c.Done():
			return false, nil
		default:
			d, err := ScanandGetDeviceByMac(out, g.MACAddress)
			if err == nil {
				return checkdeviceValueEqual(g, d)

			}
		}
	}

}

//check device info
func checkdeviceValueEqual(g gwd.NetworkConfig, d *api.DeviceInfo) (bool, []string) {
	configResults := devconfig.ConfigResult{}
	if g.NewIPAddress != d.IPAddress {
		configResults.FailFields = append(configResults.FailFields, "IPAddress")

	}
	if g.Netmask != d.Netmask {
		configResults.FailFields = append(configResults.FailFields, "Netmask")

	}
	if g.Gateway != d.Gateway {
		configResults.FailFields = append(configResults.FailFields, "Gateway")
	}

	if g.Hostname != d.Hostname {
		configResults.FailFields = append(configResults.FailFields, "Hostname")
	}
	if len(configResults.FailFields) > 0 {
		return false, configResults.FailFields
	}

	return true, nil
}

func sessionSuccess(session *devconfig.SessionState) *devconfig.SessionState {
	session.State = "success"
	session.EndedTime = time.Now().String()
	// TODO: publish session finished
	return session
}

func sessionFail(session *devconfig.SessionState, msg string) *devconfig.SessionState {
	session.State = "fail"
	session.EndedTime = time.Now().String()
	session.Message = msg
	// TODO: publish session finished
	return session
}

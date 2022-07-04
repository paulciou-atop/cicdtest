package atopudpscan

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"nms/atopudpscan/api/v1/dbstore"
	"nms/atopudpscan/pkg/gwd"
	"nms/atopudpscan/pkg/gwd/server"
	"os"
	"strconv"
	"strings"
	sync "sync"
	"time"

	"nms/atopudpscan/pkg/net"
	mq "nms/messaging"

	api "nms/api/v1/atopudpscan"

	grpc "google.golang.org/grpc"
)

var name string
var dbInf_dbstore *dbstore.DatabaseInfo

const timeout = 3

func init() {
	var err error
	name, err = os.Hostname()

	if err != nil {
		log.Fatal(err)
	}
	dbInf_dbstore = &dbstore.DatabaseInfo{
		DatabaseName: "nms_atopudpscan",
	}
}

func NewGwd() *GwdService {
	s := &GwdService{u: server.NewAtopUdpServer("0.0.0.0"), session: make([]*Session, 0), s: new(sync.Mutex)}
	mqc, _ := mq.NewClient()
	s.mq = mqc
	AtopudpscanShare().RegisterGwd(s)
	return s
}
func (g *GwdService) Run() error {
	return g.u.Run()
}

func (g *GwdService) Close() error {
	g.mq.Close()
	return g.u.Stop()
}

//create new session
func NewSession(id string) *Session {
	s := &Session{id: id, status: api.SessionStatus_running, m: new(sync.Mutex), data: make([]*api.DeviceInfo, 0), stop: false}
	return s
}

type Session struct {
	id      string
	status  api.SessionStatus
	message string
	m       *sync.Mutex
	data    []*api.DeviceInfo
	stop    bool
}

//make  seesion status to init
func (s *Session) New() {
	s.m.Lock()
	s.data = make([]*api.DeviceInfo, 0)
	s.stop = false
	s.status = api.SessionStatus_running
	s.m.Unlock()
}

//stop session
func (s *Session) Stop() {
	s.m.Lock()
	s.stop = true
	s.m.Unlock()
}

//return stop status
func (s *Session) IsStop() bool {
	s.m.Lock()
	v := s.stop
	s.m.Unlock()
	return v
}

//set session staut
func (s *Session) settingStatus(v api.SessionStatus, message string) {
	s.m.Lock()
	s.status = v
	s.message = message
	s.m.Unlock()
}

//return session staut
func (s *Session) GetStatus() api.SessionStatus {
	s.m.Lock()
	v := s.status
	s.m.Unlock()
	return v
}

//save devices to session and duplicate devices not save
func (s *Session) SettingDevice(d *api.DeviceInfo) {

	for _, v := range s.data {
		r := compare(v, d)
		if r {

			return
		}
	}
	s.data = append(s.data, d)
	log.Printf("id:%v,add vlaue:%v", s.id, d)

}

//get collection device from session
func (s *Session) GetDevice() []*api.DeviceInfo {
	v := s.data
	return v
}

type GwdService struct {
	u       *server.AtopGwdServer
	session []*Session
	s       *sync.Mutex
	api.UnimplementedGwdServer
	mq mq.IMQClient
}

//add session
func (g *GwdService) addSession(s *Session) {
	g.s.Lock()
	g.session = append(g.session, s)
	g.s.Unlock()
}

//retrun session
func (g *GwdService) getSession(id string) (*Session, error) {
	g.s.Lock()
	session := g.session
	g.s.Unlock()
	for _, v := range session {
		if v.id == id {
			return v, nil
		}
	}
	return nil, errors.New("Session not exist")
}

//make device sound
func (g *GwdService) Beep(ctx context.Context, in *api.GwdConfig) (*api.Response, error) {
	n := gwd.NetworkConfig{
		IPAddress:  in.GetIPAddress(),
		MACAddress: in.GetMACAddress(),
	}

	address := strings.Join([]string{in.GetServerIp(), strconv.Itoa(0)}, ":")
	gwd := gwd.NewAtopGwd(address)

	err := gwd.Beep(n)
	if err != nil {
		return &api.Response{Result: false, Message: err.Error()}, err
	}
	return &api.Response{Result: true}, err
}

//Stop Session scan
func (g *GwdService) Stop(ctx context.Context, in *api.Sessions) (*api.Response, error) {
	s, err := g.getSession(in.Id)
	if err != nil {
		return &api.Response{Result: false, Message: err.Error()}, err
	}
	s.Stop()
	for {
		if s.IsStop() {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	return &api.Response{Result: true}, nil
}

//return Server Ip
func (g *GwdService) GetServerIp(ctx context.Context, in *api.EmptyParams) (*api.ServerIp, error) {
	ips, err := net.GetLocalIP()
	if err != nil {
		return &api.ServerIp{}, err
	}
	return &api.ServerIp{Ip: ips}, err
}

//session scan
func (g *GwdService) SessionScan(ctx context.Context, in *api.ScanConfig) (*api.ResponseSession, error) {
	if in.Id == "" { //check id of arg is exist
		return &api.ResponseSession{Status: api.SessionStatus_fail, Message: "id is empty"}, errors.New("id is empty")
	}
	s, err := g.getSession(in.Id) //get session by id
	if err != nil {               //if not exist then create new
		s = NewSession(in.Id)
		g.addSession(s)
	} else { //if exist
		if s.IsStop() && s.GetStatus() != api.SessionStatus_running { //check if stop and not running
			s.New() //make status of session to new
		} else { //session still running
			message := fmt.Errorf("id:%v is running", in.Id)
			return &api.ResponseSession{Status: api.SessionStatus_fail, Message: message.Error()}, message
		}
	}
	address := strings.Join([]string{in.GetServerIp(), strconv.Itoa(0)}, ":")
	err = g.u.Scan(address) //scanner
	if err != nil {
		s.settingStatus(api.SessionStatus_fail, err.Error())
		s.Stop()
		return &api.ResponseSession{Status: api.SessionStatus_fail, Message: err.Error()}, err
	}
	s.settingStatus(api.SessionStatus_running, "")               //setting session status
	res := &api.ResponseSession{Id: s.id, Status: s.GetStatus()} //return value
	b := make(chan bool, 1)
	go func() { //start to collect device
		c, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		b <- true
		defer cancel()
		for {
			select {
			case <-c.Done(): //time out
				for _, v := range s.GetDevice() {
					if s.IsStop() {
						break
					}
					saveDeviceResult(s, v)
				}
				saveDeviceSession(s)                 //update session status to postgresql
				saveToDataStore(s.id, s.GetDevice()) //save collection devices of  session to datastore
				publishDone(g.mq, s.id)

				if len(s.GetDevice()) == 0 { //chech device is exist
					s.settingStatus(api.SessionStatus_notfound, "")
				} else {
					s.settingStatus(api.SessionStatus_success, "")
				}
				s.Stop() //make sure status after complete scan
				return
			default:
				d, err := g.u.GetReceiveData() //getdevice from udp
				if err == nil {
					var data []*api.DeviceInfo
					err = json.Unmarshal(d, &data)
					if err != nil {
						s.settingStatus(api.SessionStatus_fail, err.Error())
					} else {
						for _, v := range data {
							s.SettingDevice(v) //save collection device to session and postgresql
						}
						saveDeviceSession(s) //update session status to postgresql

					}
				}
				if s.IsStop() { //stop
					cancel()
				}
			}
		}

	}()
	<-b
	return res, nil

}

//Get SessionStatus
func (g *GwdService) GetSessionStatus(ctx context.Context, in *api.Sessions) (*api.ResponseSession, error) {
	s, err := g.getSession(in.GetId())
	if err != nil {
		return &api.ResponseSession{Id: in.GetId(), Status: api.SessionStatus_fail, Message: err.Error()}, err
	}

	return &api.ResponseSession{Id: s.id, Status: s.GetStatus(), Message: s.message}, nil

}

//GetSessiondata from datastore
func (g *GwdService) GetSessionData(ctx context.Context, in *api.Sessions) (*api.DeviceResponse, error) {
	dbInf_dbstore.TableName = name + "_" + in.GetId()
	h, _ := servicewatchercheck("datastore")
	/*	if err != nil {
		log.Printf("did not connect to datebase: %v", err)
		return &DeviceResponse{}, err
	}*/
	log.Printf("datastore host:%v", h)
	ct, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ct, h, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect to datebase: %v", err)
		return &api.DeviceResponse{}, fmt.Errorf("did not connect to datebase: %v", err)
	}
	c := dbstore.NewDbServiceClient(conn)
	ct, _ = context.WithTimeout(context.Background(), 3*time.Second)
	createSsRes, err := c.DbRead(ct, &dbstore.DbReadRequest{
		DbInf: dbInf_dbstore,
	})
	if err != nil {
		return &api.DeviceResponse{}, err
	}
	r := createSsRes.Rd.GetData()
	if len(r) <= 0 {
		return &api.DeviceResponse{}, errors.New("no data exist")
	}
	var deviceResponse []*api.DeviceInfo
	ma := make(map[string]interface{})
	for _, v := range r {
		m := v.Fields
		for k, v := range m {
			ma[k] = v.GetStringValue()
		}
		v := ma["IsDHCP"]
		switch v := v.(type) {
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return &api.DeviceResponse{}, err
			}
			ma["IsDHCP"] = b
		}

		data, err := json.Marshal(ma)
		if err != nil {
			return &api.DeviceResponse{}, err
		}
		var device api.DeviceInfo
		err = json.Unmarshal(data, &device)
		if err != nil {
			return &api.DeviceResponse{}, err
		}
		deviceResponse = append(deviceResponse, &device)
	}
	return &api.DeviceResponse{Devices: deviceResponse}, nil
}

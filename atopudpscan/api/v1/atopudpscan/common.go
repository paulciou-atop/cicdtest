package atopudpscan

import sync "sync"

var once sync.Once
var s *ServiceComm

func AtopudpscanShare() *ServiceComm {
	once.Do(func() {
		s = &ServiceComm{l: new(sync.Mutex)}
	})

	return s
}

//for share others service
type ServiceComm struct {
	commonGwd              *GwdService
	commonDeviceController *DeviceController
	l                      *sync.Mutex
}

func (s *ServiceComm) RegisterGwd(g *GwdService) {
	s.l.Lock()
	s.commonGwd = g
	s.l.Unlock()
}

func (s *ServiceComm) GetGwd() *GwdService {
	s.l.Lock()
	g := s.commonGwd
	s.l.Unlock()
	return g
}

func (s *ServiceComm) RegisterDeviceController(d *DeviceController) {
	s.l.Lock()
	s.commonDeviceController = d
	s.l.Unlock()
}

func (s *ServiceComm) GetDeviceController() *DeviceController {
	s.l.Lock()
	d := s.commonDeviceController
	s.l.Unlock()
	return d
}

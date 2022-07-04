package devconfig

import (
	"context"
	"errors"
	"fmt"
	api "nms/api/v1/devconfig"
	"sync"

	"nms/config/pkg/config"
	"nms/config/pkg/session"
	"nms/lib/repo"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var ErrNotImplement = errors.New("Not implmented")

type Server struct {
	api.UnimplementedConfigServer
}

func RegisterServices(s *grpc.Server) error {
	api.RegisterConfigServer(s, &Server{})
	return nil
}

// Upload upload configuration from specific device
func (s *Server) Upload(ctx context.Context, req *api.UploadConfigRequest) (*api.UploadConfigResponse, error) {
	if req.Device == nil {
		return nil, fmt.Errorf("null device")
	}
	ctx, cancel := context.WithTimeout(context.Background(), CONIFG_TIMEOUT)
	defer cancel()
	r, err := repo.GetRepo(ctx)
	if err != nil {
		return nil, err
	}
	c := config.NewConfig(r)
	dev := config.Device{
		ID:   req.Device.DeviceId,
		Path: req.Device.DevicePath,
	}
	ret, err := c.Upload(ctx, dev, req.Protol, req.Kinds)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Config service
//Device config specific device Config.Device
func (s *Server) Device(request *api.DeviceConfigRequest, stream api.Config_DeviceServer) error {

	if request.Device == nil {
		return fmt.Errorf("null device")
	}

	if request.Settings == nil {
		return fmt.Errorf("null settings")
	}

	device := config.Device{
		ID:   request.Device.DeviceId,
		Path: request.Device.DevicePath,
	}

	metrics := config.MarshalConfigMetrics(request.Settings)
	ctx, cancel := context.WithTimeout(context.Background(), CONIFG_TIMEOUT)
	defer cancel()
	r, err := repo.GetRepo(ctx)
	if err != nil {
		return err
	}
	c := config.NewConfig(r)
	done := make(chan int)
	res, err := c.Config(ctx, device, metrics, request.Valid, done)
	if err != nil {
		return err
	}

	stream.Send(res)
	//TODO - subscribe session done and send result
	<-done

	return nil
}

var CONIFG_TIMEOUT = config.CONFIG_TIMEOUT

// Devices config multiple devices
func (s *Server) Devices(request *api.DevicesConfigRequest, stream api.Config_DevicesServer) error {
	var devices = []config.Device{}
	for _, dev := range request.Devices {
		device := config.Device{
			ID:   dev.DeviceId,
			Path: dev.DevicePath,
		}
		devices = append(devices, device)
	}
	metrics := config.MarshalConfigMetrics(request.Settings)
	ctx, cancel := context.WithTimeout(context.Background(), CONIFG_TIMEOUT)
	defer cancel()

	r, err := repo.GetRepo(ctx)
	if err != nil {
		return err
	}
	c := config.NewConfig(r)
	wg := new(sync.WaitGroup)
	wg.Add(len(devices))
	for _, dev := range devices {
		done := make(chan int)
		go func(dev config.Device) {
			res, err := c.Config(ctx, dev, metrics, request.Valid, done)
			if err != nil {
				logrus.Errorf("config dev %s fail %v", dev.ID, err)
			}

			stream.Send(res)
			<-done
			wg.Done()
		}(dev)

	}
	wg.Wait()
	return nil
}

// Store device's configuration file
func (s *Server) Store(ctx context.Context, request *api.StoreRequest) (*api.StoreResponse, error) {
	return nil, ErrNotImplement
}

func (s *Server) GetResult(ctx context.Context, request *api.GetResultRequest) (*api.GetResultResponse, error) {
	r, err := repo.GetRepo(ctx)
	if err != nil {
		return nil, err
	}
	db := r.DB()
	configSession, err := session.GetConfigSession(db, request.SessionId)
	if err != nil {
		return nil, err
	}
	c := config.NewConfig(r)
	sessionState := session.UnMarshalSessionState(&configSession)
	configResults, err := c.GetConfigResults(request.SessionId)
	if err != nil {
		return nil, err
	}
	return &api.GetResultResponse{
		Session:       sessionState,
		ConfigResults: configResults,
	}, nil
}

func (s *Server) List(ctx context.Context, request *api.ListRequest) (*api.ListResponse, error) {
	r, err := repo.GetRepo(ctx)
	if err != nil {
		return nil, err
	}

	db := r.DB()

	var sessions = []session.ConfigSession{}

	if request.SessionIds == nil || len(request.SessionIds) <= 0 {
		//Get All
		sessions, err = session.GetAllConfigSession(ctx, db)
		if err != nil {
			return nil, err
		}
	} else {
		for _, id := range request.SessionIds {
			s, err := session.GetConfigSession(db, id)
			if err != nil {
				return nil, err
			}
			sessions = append(sessions, s)
		}
	}
	var states = []*api.SessionState{}
	for _, s := range sessions {
		state := session.UnMarshalSessionState(&s)
		states = append(states, state)
	}
	return &api.ListResponse{Sessions: states}, nil
}

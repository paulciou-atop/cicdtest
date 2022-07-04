package serviceswatcher

import (
	"context"
	"errors"
	"fmt"
	api "nms/api/v1/serviceswatcher"

	"nms/serviceswatcher/pkg/watcher"

	"google.golang.org/grpc"
)

var errNotImplement = errors.New("not implement")
var errNotFound = func(service string) error {
	return fmt.Errorf("service %s did not find", service)
}

func RegisterServices(s *grpc.Server) error {
	api.RegisterWatcherServer(s, &server{})
	return nil
}

type server struct {
	api.UnimplementedWatcherServer
}

// List API handler
func (s *server) List(ctx context.Context, req *api.Empty) (*api.ServicesResponse, error) {
	ret := &api.ServicesResponse{}
	serviceList := watcher.ListServices()
	for _, v := range serviceList {
		i := &api.ServiceInfo{
			Name:    v.Name,
			Address: v.Address,
			Port:    v.Port,
			Kind:    v.Kind,
		}
		ret.Infos = append(ret.Infos, i)
	}

	return ret, nil

}

//Get API handler
func (s *server) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	info, err := watcher.GetService(req.ServiceName)
	if err != nil {
		return &api.GetResponse{Success: false}, errNotFound(req.ServiceName)
	}
	return &api.GetResponse{
		Success: true,
		Info: &api.ServiceInfo{
			Name:    info.Name,
			Address: info.Address,
			Port:    info.Port,
			Kind:    info.Kind,
		},
	}, nil
}

//Register API handler
func (s *server) Register(ctx context.Context, req *api.ServiceInfo) (*api.RegisterResponse, error) {
	err := watcher.RegisterService(watcher.ServiceInfo{
		Name:    req.Name,
		Address: req.Address,
		Port:    req.Port,
		Kind:    req.Kind,
	})
	if err != nil {
		return &api.RegisterResponse{Success: false, Reason: err.Error()}, err
	}
	return &api.RegisterResponse{Success: true}, nil
}

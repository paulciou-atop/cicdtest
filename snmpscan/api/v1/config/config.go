package config

import (
	"context"
	"errors"
	api "nms/api/v1/config"
	"nms/snmpscan/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

var ErrNotImplement = errors.New("not implement")

func RegisterServices(s *grpc.Server) error {
	api.RegisterConfigServer(s, &service{})
	return nil
}

type service struct {
	api.UnimplementedConfigServer
}

func (s *service) Set(ctx context.Context, req *api.ConfigRequest) (*api.ConfigResponse, error) {
	preConfig := viper.AllSettings()
	updatedConfig := config.UpdateConfig(req.Config.AsMap())

	pre, err := structpb.NewStruct(preConfig)
	if err != nil {
		log.Error("config.Set convert map[string]interface{} to pb.Struct err: ", err)
		return &api.ConfigResponse{}, err
	}
	updated, err := structpb.NewStruct(updatedConfig)
	if err != nil {
		log.Error("config.Set convert map[string]interface{} to pb.Struct err: ", err)
		return &api.ConfigResponse{}, err
	}
	return &api.ConfigResponse{
		PreConfig: pre,
		NewConfig: updated,
	}, nil
}

// Update update configuration
func (s *service) Update(c context.Context, r *api.ConfigRequest) (*api.ConfigResponse, error) {
	return s.Set(c, r)
}

// Get get all conifguration, return {preConfig, newConfig} preConfig and newConfig are same data
func (s *service) Get(context.Context, *api.GetRequest) (*api.ConfigResponse, error) {
	settings := viper.AllSettings()
	ret, err := structpb.NewStruct(settings)
	if err != nil {
		log.Error("config.Set convert map[string]interface{} to pb.Struct err: ", err)
		return &api.ConfigResponse{}, err
	}
	return &api.ConfigResponse{
		PreConfig: ret,
		NewConfig: ret,
	}, nil
}

/*
devconfig.proto client code
*/
package configer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	api "nms/api/v1/configer"
	devconfig "nms/api/v1/devconfig"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type Server struct {
	api.UnimplementedConfigerServer
}

// fileConfig config process
func fileConfig(path string, config *devconfig.ConfigOptions) (*devconfig.ConfigResult, error) {
	jsonFile, err := os.Open(path)
	allFail := allFail(config.Payload.AsMap())
	if err != nil {

		return &devconfig.ConfigResult{
			Protocol:   config.Protocol,
			Kind:       config.Kind,
			Hash:       config.Hash,
			FailFields: allFail,
		}, fmt.Errorf("file %s not found", path)
	}

	// check if any unsupport option
	unsupported := testUnsupport(config.Kind, config.Payload.AsMap())
	if len(unsupported) > 0 {
		return &devconfig.ConfigResult{
			Protocol:   config.Protocol,
			Kind:       config.Kind,
			Hash:       config.Hash,
			FailFields: unsupported,
		}, fmt.Errorf("had unsupported options %v", unsupported)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	var deviceConfg map[string]interface{}
	json.Unmarshal([]byte(byteValue), &deviceConfg)
	deviceConfg[config.Kind] = config.Payload.AsMap()

	updatedContent, _ := json.MarshalIndent(deviceConfg, "", "  ")
	err = ioutil.WriteFile(path, updatedContent, 0644)
	if err != nil {
		return &devconfig.ConfigResult{
			Protocol:   config.Protocol,
			Kind:       config.Kind,
			Hash:       config.Hash,
			FailFields: allFail,
		}, fmt.Errorf("write to file fail: %v", err)
	}
	//success
	return &devconfig.ConfigResult{
		Protocol:   config.Protocol,
		Kind:       config.Kind,
		Hash:       config.Hash,
		FailFields: []string{},
	}, nil
}

func checkDeviceExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func (s *Server) Config(stream api.Configer_ConfigServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
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
		// configer should has config function instead of anonymous closure function

		// 1. Check device exist
		if !checkDeviceExist(req.Device.DevicePath) {
			logrus.Infof("device %s not exist", req.Device.DevicePath)
			//file not found whole config fail
			session := req.Session
			session = sessionFail(session, fmt.Sprintf("can not reach device %s", req.Device.DevicePath))
			var configResults = []*devconfig.ConfigResult{}
			// Get each config
			for _, c := range req.Configs {
				failFields := allFail(c.Payload.AsMap())
				failconfig := devconfig.ConfigResult{
					Protocol:   c.Protocol,
					Kind:       c.Kind,
					Hash:       c.Hash,
					FailFields: failFields,
				}
				configResults = append(configResults, &failconfig)
			}
			err = stream.Send(&api.ConfigerResponse{
				Session:       session,
				Device:        req.Device,
				ConfigResults: configResults,
			})
			if err != nil {
				return err
			}
			continue

		}
		// 2. device found, then configure it
		var configResults = []*devconfig.ConfigResult{}
		session := req.Session
		success := true
		for _, config := range req.Configs {
			// do config
			configResult, err := fileConfig(req.Device.DevicePath, config)
			if len(configResult.FailFields) > 0 {
				success = false
			}
			configResults = append(configResults, configResult)
			if err != nil {
				logrus.Error("config err ", err)
			}

		}

		if !success {
			logrus.Infof("config %s fail", session.Id)
			session = sessionFail(session, "some fields fail")
		} else {
			logrus.Infof("config %s success", session.Id)
			session = sessionSuccess(session)
		}
		err = stream.Send(&api.ConfigerResponse{
			Session:       session,
			Device:        req.Device,
			ConfigResults: configResults,
		})
		if err != nil {
			return err
		}

	}
	return nil
}
func (s *Server) Validate(stream api.Configer_ValidateServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		var configResults = []*devconfig.ConfigResult{}
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

		for _, config := range configs {
			unsupported := testUnsupport(config.Kind, config.Payload.AsMap())

			if len(unsupported) > 0 {
				configResults = append(configResults, &devconfig.ConfigResult{
					Protocol:   config.Protocol,
					Kind:       config.Kind,
					Hash:       config.Hash,
					FailFields: unsupported,
				})
			}
		}

		if len(configResults) > 0 {
			session = sessionFail(session, fmt.Sprintf("some config not support"))
		} else {
			session = sessionSuccess(session)
		}

		stream.Send(&api.ValidateResponse{
			Session:       session,
			ConfigResults: configResults,
		})

	}
	return nil
}
func (s *Server) FileTransfer(stream api.Configer_FileTransferServer) error {
	return fmt.Errorf("not support")
}

func (s *Server) GetConfig(ctx context.Context, req *api.GetConfigRequest) (*api.GetConfigResponse, error) {
	if !checkDeviceExist(req.Device.DevicePath) {
		return &api.GetConfigResponse{}, fmt.Errorf("device not found")
	}
	jsonFile, err := os.Open(req.Device.DevicePath)
	if err != nil {
		return &api.GetConfigResponse{}, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return &api.GetConfigResponse{}, err
	}

	var deviceConfg map[string]interface{}
	err = json.Unmarshal(byteValue, &deviceConfg)
	if err != nil {
		return &api.GetConfigResponse{}, err
	}

	var responseConfig = map[string]interface{}{}
	if len(req.Kinds) <= 0 {
		//All
		responseConfig = deviceConfg
	} else {
		for _, k := range req.Kinds {
			responseConfig[k] = deviceConfg[k]
		}
	}

	m, err := structpb.NewStruct(responseConfig)

	if err != nil {
		return &api.GetConfigResponse{}, err
	}

	return &api.GetConfigResponse{
		Device:  req.Device,
		Configs: m,
	}, nil
}

func RegisterServices(s *grpc.Server) error {
	api.RegisterConfigerServer(s, &Server{})
	return nil
}

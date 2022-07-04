package config

import (
	"context"
	"fmt"
	"io"
	"time"

	"nms/api/v1/common"
	configerAPI "nms/api/v1/configer"
	"nms/api/v1/devconfig"

	"nms/config/internal/services"
	"nms/config/pkg/session"
	"nms/lib/pgutils"
	"nms/lib/repo"

	"github.com/sirupsen/logrus"
)

type IConfig interface {
	Config(ctx context.Context, dev Device, metrics []*ConfigMetric, valid bool, done chan<- int) (*devconfig.ConfigResponse, error)
	GetConfigResults(sessionid string) ([]*devconfig.ConfigResult, error)
	Upload(ctx context.Context, dev Device, protocol string, kinds []string) (*devconfig.UploadConfigResponse, error)
}

func NewConfig(parent repo.IRepo) IConfig {
	return &configWithRepo{r: parent}
}

type configWithRepo struct {
	r repo.IRepo
}

func (cr *configWithRepo) validate(ctx context.Context, c chan<- configerAPI.ConfigerClient, protocol string, metrics []*ConfigMetric) {

	conn := cr.r.Connection(protocol)
	if conn == nil {
		return
	}
	client := configerAPI.NewConfigerClient(conn)

	stream, err := client.Validate(ctx)
	if err != nil {
		logrus.Errorf("%s create stream fail %v", protocol, err)
		return
	}
	go func() {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				// we don't care error, we just want to find any service can config
				return
			}
			failCount := 0
			for _, res := range result.ConfigResults {
				failCount += len(res.FailFields)
			}

			if failCount == 0 {
				c <- client
				return
			}
		}
	}()

	configs, err := UnMarshalAPIConfigOptions(metrics)
	if err != nil {
		logrus.Errorf("convert config options fail: %v", err)
		return
	}
	// We don't care session in this case, because we only want to find
	// a configer accept metrics
	stream.Send(&configerAPI.ConfigerValidateRequest{
		Session: &devconfig.SessionState{
			Id:          "doesn't matter",
			State:       "running",
			StartedTime: time.Now().String(),
		},
		Configs: configs,
	})
}

type Device struct {
	ID   string
	Path string
}

// TODO timeout should from config
var TIMEOUT = time.Second * 5

// getConfigerClient find suitable configer client
func (cr *configWithRepo) getConfigerClient(ctx context.Context, metrics []*ConfigMetric, valid bool) (configerAPI.ConfigerClient, error) {

	out := make(chan configerAPI.ConfigerClient)
	serviceNames := services.ConfigerList
	go func() {
		if valid {
			// get first configer
			for _, clientName := range serviceNames {
				go cr.validate(ctx, out, clientName, metrics)
			}

		} else {
			clientName := metrics[0].protocol

			conn := cr.r.Connection(clientName)

			if conn != nil {
				client := configerAPI.NewConfigerClient(conn)
				//success
				out <- client
			}
			// let it timeout if error
		}
	}()

	timeout := time.After(TIMEOUT)
	var client configerAPI.ConfigerClient
	select {
	case client = <-out:

	case <-ctx.Done():
		logrus.Info("cancel config")
		return nil, ctx.Err()
	case <-timeout:
		logrus.Error("valid timeout ")
		return nil, fmt.Errorf("find configer timeout")
	}

	return client, nil

}

// Upload
func (c *configWithRepo) Upload(ctx context.Context, device Device, protocol string, kinds []string) (*devconfig.UploadConfigResponse, error) {
	conn := c.r.Connection(protocol)
	if conn == nil {
		return nil, fmt.Errorf("can not find protocol: %s", protocol)
	}
	client := configerAPI.NewConfigerClient(conn)

	dev := &common.DeviceIdentify{
		DeviceId:   device.ID,
		DevicePath: device.Path,
	}
	ret, err := client.GetConfig(ctx, &configerAPI.GetConfigRequest{
		Device: dev,
		Kinds:  kinds,
	})
	if err != nil {
		return nil, err
	}
	res := &devconfig.UploadConfigResponse{
		Success: true,
		Message: "",
		Device:  dev,
		Protol:  protocol,
		Kinds:   kinds,
		Payload: ret.Configs,
	}
	return res, nil
}

// Config configure specific device
func (c *configWithRepo) Config(ctx context.Context, device Device, metrics []*ConfigMetric, valid bool, done chan<- int) (*devconfig.ConfigResponse, error) {
	//validate conifg for each configer, process configuration on the first configer accept metrics
	// receive validate stream
	if len(metrics) <= 0 {
		return nil, fmt.Errorf("empty config metrics")
	}

	client, err := c.getConfigerClient(ctx, metrics, valid)
	if err != nil {
		return nil, err
	}
	// get store interface <-- side effect here
	store := c.r.DB()

	if err := storeMetrics(store, metrics); err != nil {
		return nil, err
	}

	// new session
	sChan := session.NewSession(ctx)
	s := <-sChan
	s.WriteToStore(store)

	// create configer stream
	stream, err := client.Config(ctx)
	if err != nil {
		logrus.Errorf("configer create stream fail :%v ", err)
		return nil, err
	}

	chanNewTasks := make(chan []session.ConfigTask)
	chanUpdateSession := make(chan session.ConfigSession)
	// response routine
	go func() {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			logrus.Errorf("receive config response fail: %v", err)
			return
		}

		resTasks, err := session.MarshalConfigTasks(res)
		if err != nil {
			logrus.Errorf("marshal config tasks fail: %v", err)
			return
		}

		sc, err := session.MarshalConfigSession(res)
		if err != nil {
			logrus.Errorf("marshal config response fail: %v", err)
			return
		}
		chanUpdateSession <- sc
		chanNewTasks <- resTasks
		return
	}()

	options, err := UnMarshalAPIConfigOptions(metrics)
	if err != nil {
		logrus.Errorf("convert config option fail: %v", err)
		return nil, err
	}

	// send request
	stream.Send(&configerAPI.ConfigerConfigRequest{
		Session: session.UnMarshalSessionState(&s.State),
		Device: &common.DeviceIdentify{
			DeviceId:   device.ID,
			DevicePath: device.Path,
		},
		Configs: options,
	})
	stream.CloseSend()

	// got response or timeout
	go func() {
		timeout := time.After(TIMEOUT)
		for {
			select {
			case sc := <-chanUpdateSession:
				s.Update(sc)
			case newTasks := <-chanNewTasks:
				s.AddTasks(newTasks)
				s.Done()
				goto WriteBack
			case <-timeout:
				s.State.Cancel()
				logrus.Error("waiting configer response timeout ")
				goto WriteBack
			}
		}
	WriteBack:
		s.WriteToStore(store)
		s.Publish(c.r.MQ())
		done <- 1
	}()

	ret := &devconfig.ConfigResponse{
		Session: session.UnMarshalSessionState(&s.State),
		Device: &common.DeviceIdentify{
			DeviceId:   device.ID,
			DevicePath: device.Path,
		},
	}
	return ret, nil

}

// GetConfigResults get devconfig.ConfigResult
func (c *configWithRepo) GetConfigResults(sessionid string) ([]*devconfig.ConfigResult, error) {
	querier := c.r.DB()
	var results []*devconfig.ConfigResult
	tasks, err := session.GetConfigTasks(querier, sessionid)
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		var m ConfigMetricModule
		querier.Query(&m, pgutils.QueryExpr{
			Expr:  "hash = ?",
			Value: t.ConfigHash,
		})
		r := &devconfig.ConfigResult{
			Protocol:   m.Protocol,
			Kind:       m.Kind,
			Hash:       m.Hash,
			FailFields: t.FaildOptions,
		}
		results = append(results, r)
	}
	return results, nil
}

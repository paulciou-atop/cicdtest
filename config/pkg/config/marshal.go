package config

import (
	"nms/api/v1/devconfig"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
)

func UnMarshalAPIConfigOptions(metrics []*ConfigMetric) ([]*devconfig.ConfigOptions, error) {
	var configs = []*devconfig.ConfigOptions{}

	for _, m := range metrics {
		payload, err := structpb.NewStruct(m.payload)
		if err != nil {
			logrus.Errorf("create config option fail: %v", err)
			return nil, err
		}
		configs = append(configs, &devconfig.ConfigOptions{
			Protocol: m.protocol,
			Kind:     m.kind,
			Hash:     m.hash,
			Payload:  payload,
		})
	}
	return configs, nil
}

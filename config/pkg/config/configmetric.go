/*
ConfigMetricModule use when we want to create a database table.
ConfigMetric defined configuration and related information.
*/
package config

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"nms/api/v1/devconfig"
)

type Metric = map[string]interface{}

// ConfigMetricModule use this struct to create postgre table
type ConfigMetricModule struct {
	tableName  struct{} `pg:"config_metrics,alias:conf"`
	ID         int64
	Name       string
	Protocol   string
	Kind       string
	Payload    Metric
	Hash       string
	Count      int32
	LastConfig string
}

// ConfigMetric
type ConfigMetric struct {
	protocol string
	kind     string
	payload  Metric
	hash     string
}

func (c ConfigMetric) Hash() string {
	return c.hash
}
func (c ConfigMetric) Payload() Metric {
	return c.payload
}

func hash(s Metric) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(s)
	h := sha256.New()
	h.Write(b.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func NewConfigMetric(protocol string, kind string, m Metric) *ConfigMetric {
	return &ConfigMetric{
		protocol: protocol,
		kind:     kind,
		payload:  m,
		hash:     hash(m),
	}
}

func MarshalConfigMetrics(opts []*devconfig.ConfigOptions) []*ConfigMetric {
	var metrics = []*ConfigMetric{}
	for _, o := range opts {
		m := NewConfigMetric(o.Protocol, o.Kind, o.Payload.AsMap())
		metrics = append(metrics, m)
	}
	return metrics
}

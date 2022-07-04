package config

import (
	"nms/config/pkg/session"
	"nms/lib/pgutils"
	pg "nms/lib/pgutils"
	"time"

	"github.com/sirupsen/logrus"
)

func InitTable(client pg.IClient) error {
	err := client.CreateTable(&ConfigMetricModule{}, pg.CreateTableOpt{IfNotExists: true})
	if err != nil {
		logrus.Error("create config metric state table fail: ", err)

	}
	return nil
}

func storeMetrics(db pgutils.IDataAccess, metrics []*ConfigMetric) error {
	for _, v := range metrics {
		if err := StoreConfig(db, v); err != nil {
			return err
		}
	}
	return nil
}

// StoreConfig store ConfigMetric into database
func StoreConfig(db pgutils.IDataAccess, config *ConfigMetric) error {
	var configs []ConfigMetricModule
	// check config is exist by hash
	err := db.Query(&configs, pgutils.QueryExpr{
		Expr:  "hash = ?",
		Value: config.hash,
	})
	if err != nil {
		logrus.Error("query config metric fail: ", err)
		return err
	}
	// not found insert new
	if len(configs) < 1 {
		if err := db.Insert(&ConfigMetricModule{
			Protocol:   config.protocol,
			Kind:       config.kind,
			Payload:    config.payload,
			Hash:       config.hash,
			LastConfig: time.Now().Format(session.TIME_LAYOUT),
			Count:      1,
		}); err != nil {
			logrus.Error("store config metric fail: ", err)
			return err
		}
		return nil
	}
	// config exist count+1
	confAddCount := configs[0]
	confAddCount.Count += 1
	confAddCount.LastConfig = time.Now().Format(session.TIME_LAYOUT)
	if err := db.Update(&confAddCount); err != nil {
		logrus.Error("store config metric fail: ", err)
		return err
	}
	return nil
}

package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Object = map[string]interface{}

//setDefaultSettings set default setting if configuration do not exist
func setDefaultSettings() {
	viper.SetDefault("services", []Object{
		{"name": "servicewatcher",
			"address": "servicewatcher",
			"port":    8081,
			"kind":    []string{"https", "grpc"}},
		{
			"name":    "scanservice",
			"address": "scanservice",
			"port":    8080,
			"kind":    []string{"https", "grpc"},
		},
		{
			"name":    "snmpscan",
			"address": "snmpscan",
			"port":    8084,
			"kind":    []string{"https", "grpc"},
		},
		{
			"name":    "redis",
			"address": "redis",
			"port":    6379,
			"kind":    []string{"redis"},
		},
		{
			"name":    "atopudpscan",
			"address": "atopudpscan",
			"port":    8083,
			"kind":    []string{"https", "grpc"},
		},
		{
			"name":    "mongo",
			"address": "mongo",
			"port":    27017,
			"kind":    []string{"db"},
		},
		{
			"name":    "datastore",
			"address": "datastore",
			"port":    8082,
			"kind":    []string{"db"},
		},
		{
			"name":    "dummyconfiger",
			"address": "dummyconfiger",
			"port":    8085,
			"kind":    []string{"grpc"},
		},
		{
			"name":    "config",
			"address": "config",
			"port":    8100,
			"kind":    []string{"grpc"},
		},
		{
			"name":    "inventory",
			"address": "inventory",
			"port":    8101,
			"kind":    []string{"grpc"},
		},
	})

	viper.SetDefault("mode", "container")
}

func init() {
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	setDefaultSettings()
	if err != nil {
		logrus.Warning("Configuration file not found, use default instead of: ", err)
		logrus.Infoln("save default settings into settings.json")
	}
	if err = viper.WriteConfigAs("./settings.json"); err != nil {
		logrus.Error("write configuration file fail : ", err)
	}
}

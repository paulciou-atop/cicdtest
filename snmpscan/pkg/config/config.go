/*
Package config implements bunch of function about configuration. Configuration usually means
default settings of this service.

*/
package config

import (
	glog "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Object = map[string]interface{}

//setDefaultSettings set default setting if configuration do not exist
func setDefaultSettings() {
	viper.SetDefault("snmp", Object{
		"port":      161,
		"community": "public",
		"version":   0x01,
		"timeout":   1,
		"Retries":   1,
	})

	viper.SetDefault("atopDevices", []string{
		".1.3.6.1.4.1.3755.0.0.14",   //  eh7506
		".1.3.6.1.4.1.3755.0.0.23",   // eh7508
		".1.3.6.1.4.1.3755.0.0.15",   //  eh7512
		".1.3.6.1.4.1.3755.0.0.21",   //  eh7520, ehg750x
		".1.3.6.1.4.1.3755.0.0.22",   //  rh7728
		".1.3.6.1.4.1.3755.0.0.24",   //  eh1504, eh7504
		".1.3.6.1.4.1.3755.0.0.26",   //  rh7528
		".1.3.6.1.4.1.3755.0.0.31",   //  ehg760x L2
		".1.3.6.1.4.1.3755.0.0.33",   //  ehg760x L3
		".1.3.6.1.4.1.3755.0.0.2033", //  ehg760x L3
		".1.3.6.1.4.1.3755.0.0.2031", //  ehg760x L3
		".1.3.6.1.4.1.3755.0.0.2015", //  ehg760x L3
		".1.3.6.1.4.1.11317.0.0.31",  // Wieland L2MS 4G/8G
		".1.3.6.1.4.1.11317.0.0.14",  // Wieland L2MS 06
		".1.3.6.1.4.1.11317.0.0.23",  // Wieland L2MS 08
		".1.3.6.1.4.1.11317.0.0.15",  // Wieland L2MS 12
		".1.3.6.1.4.1.11317.0.0.21",  // Wieland L2MS 20
	},
	)
}

func init() {
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	setDefaultSettings()
	if err != nil {
		glog.Warning("Configuration file not found, use default instead of: ", err)
		glog.Infoln("save default settings into settings.json")
	}
	if err = viper.WriteConfigAs("./settings.json"); err != nil {
		glog.Error("write configuration file fail : ", err)
	}
}

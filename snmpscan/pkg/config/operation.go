package config

import (
	"github.com/spf13/viper"
)

//UpdateConfig
func UpdateConfig(updateConfigs Object) Object {
	for k, v := range updateConfigs {
		viper.Set(k, v)
	}
	viper.WriteConfig()
	return viper.AllSettings()
}

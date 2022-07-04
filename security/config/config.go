package config

import (
	"log"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(getCurrentAbPathByCaller())
	err := viper.ReadInConfig()
	if err != nil {
		setDefault()
		viper.WriteConfig()
		if err = viper.WriteConfigAs(getCurrentAbPathByCaller() + "/config.json"); err != nil {
			log.Println("can not written setting file.")
		}
		log.Println("writting default config file")
	}
}

func setDefault() {
	grpcConfigInit()
	pkiConfigInit()
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

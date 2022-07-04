package configs

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

const grpc = "grpcPort"
const http = "httpPort"

const grpcPort = 8083
const httpPort = 8093

func GetgrpcPort() string {
	port := viper.GetString(grpc)

	return port
}
func GethttpPort() string {
	port := viper.GetString(http)

	return port
}

func setDefault() {
	viper.SetDefault(grpc, grpcPort)
	viper.SetDefault(http, httpPort)

}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

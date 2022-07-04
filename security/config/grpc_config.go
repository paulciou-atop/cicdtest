package config

import "github.com/spf13/viper"

const grpc = "grpcPort"
const http = "httpPort"
const grpcPort = 8080
const httpPort = 8090

func grpcConfigInit() {
	viper.SetDefault(grpc, grpcPort)
	viper.SetDefault(http, httpPort)
}

func GetgrpcPort() string {
	port := viper.GetString(grpc)

	return port
}

func GethttpPort() string {
	port := viper.GetString(http)

	return port

}

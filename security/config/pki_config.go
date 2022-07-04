package config

import (
	"net"

	"github.com/spf13/viper"
)

const keyname = "Name"
const valuename = "Atop Inc."

const keydns = "Dns"
const valuedns = "localhost"

const keyaddress = "Address"
const valueaddress = "127.0.0.1"

const keyprovisioner = "Provisioner"
const valueprovisioner = "Atop@example.com"

const keypwdFileName = "PwdFileName"
const valuepwdFileName = "password"

const keypassword = "password"
const valuepassword = "atop"

const keyport = "Port"
const valueport = 8443

const token = "token"

func pkiConfigInit() {
	viper.SetDefault(keyname, valuename)
	viper.SetDefault(keydns, valuedns)
	viper.SetDefault(keyaddress, valueaddress)
	viper.SetDefault(keyprovisioner, valueprovisioner)
	viper.SetDefault(keypwdFileName, valuepwdFileName)
	viper.SetDefault(keypassword, valuepassword)
	viper.SetDefault(keyport, valueport)
}

func SetToken(t string) {
	viper.Set(token, t)
	viper.WriteConfig()
}

func GetToken() string {
	return viper.GetString(token)
}

func GetpwdFileName() string {
	return viper.GetString(keypwdFileName)
}

func Getpassword() string {
	return viper.GetString(keypassword)
}

func GetDns() string {
	return viper.GetString(keydns)
}

func getPort() string {
	return viper.GetString(keyport)
}

func GetAddress() string {
	return net.JoinHostPort(viper.GetString(keyaddress), getPort())
}

func GetCaUrl() string {
	return net.JoinHostPort(viper.GetString(keydns), getPort())
}

func GetProvisioner() string {
	return viper.GetString(keyprovisioner)
}

func GetName() string {
	return viper.GetString(keyname)
}

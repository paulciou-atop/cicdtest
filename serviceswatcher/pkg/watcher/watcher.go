package watcher

import (
	"fmt"
	_ "nms/serviceswatcher/config"
	"os"

	"github.com/spf13/viper"
)

type Object = map[string]interface{}

var ErrNotFound = func(name string) error {
	return fmt.Errorf("service %s did not found", name)
}

type ServiceInfo struct {
	Name    string   `json:"name"`
	Address string   `bson:"address" json:"address" structs:"address"`
	Port    int32    `bson:"port" json:"port" structs:"port"`
	Kind    []string `bson:"kind" json:"kind" structs:"kind"`
}

func isRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

func host(info ServiceInfo) string {
	h := info.Address
	if !isRunningInDockerContainer() {
		h = "localhost"
	}
	return fmt.Sprintf("%s:%d", h, info.Port)
}

// ListServices list all servicews
func ListServices() []ServiceInfo {
	ret := []ServiceInfo{}
	if err := viper.GetViper().UnmarshalKey("services", &ret); err != nil {
		return []ServiceInfo{}
	}
	return ret
}

// GetService get service information
func GetService(name string) (ServiceInfo, error) {

	services := []ServiceInfo{}
	if err := viper.GetViper().UnmarshalKey("services", &services); err != nil {
		return ServiceInfo{}, err
	}
	for _, v := range services {
		if v.Name == name {
			return v, nil
		}
	}
	return ServiceInfo{}, ErrNotFound(name)
}

func GetServiceHostUrl(name string) (string, error) {
	services := []ServiceInfo{}
	if err := viper.GetViper().UnmarshalKey("services", &services); err != nil {
		return "", err
	}
	for _, v := range services {
		if v.Name == name {
			return host(v), nil
		}
	}
	return "", ErrNotFound(name)
}

// RegisterService register new service or update exist service
func RegisterService(newService ServiceInfo) error {
	services := []ServiceInfo{}
	if err := viper.GetViper().UnmarshalKey("services", &services); err != nil {
		return err
	}
	for i, v := range services {
		if v.Name == newService.Name {
			services[i] = ServiceInfo{
				Name:    newService.Name,
				Address: newService.Address,
				Port:    newService.Port,
				Kind:    newService.Kind,
			}
			viper.Set("services", services)
			viper.WriteConfig()
			return nil
		}
	}
	services = append(services, ServiceInfo{
		Name:    newService.Name,
		Address: newService.Address,
		Port:    newService.Port,
		Kind:    newService.Kind,
	})
	viper.Set("services", services)
	viper.WriteConfig()
	return nil
}

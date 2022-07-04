/*
Copyright © 2022 Atop NMS team

*/
package register

import (
	"fmt"

	"nms/serviceswatcher/pkg/watcher"

	glog "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RegsiterOption struct {
	Name    string //service's name
	Address string
	Port    int32
	Kind    []string
}

// NewCmdRegister register specific service information
func NewCmdRegister() *cobra.Command {
	o := &RegsiterOption{}
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "register a new service",
		Run: func(cmd *cobra.Command, args []string) {
			err := watcher.RegisterService(watcher.ServiceInfo{
				Name:    o.Name,
				Address: o.Address,
				Port:    o.Port,
				Kind:    o.Kind,
			})
			if err != nil {
				glog.Errorf("Register service error %+v", err)
				return
			}

			fmt.Printf("service %s was registered on %s:%d ", o.Name, o.Address, o.Port)
		}}
	registerCmd.Flags().StringArrayVarP(&o.Kind, "kind", "k", []string{}, "What kind of API supported")
	registerCmd.Flags().StringVarP(&o.Name, "name", "n", "unknown", "service's name")
	registerCmd.Flags().StringVar(&o.Address, "address", "localhost", "service's hostname/IP address")
	registerCmd.Flags().Int32VarP(&o.Port, "port", "p", 8080, "service's port number")

	registerCmd.MarkFlagRequired("name")
	registerCmd.MarkFlagRequired("address")
	registerCmd.MarkFlagRequired("port")
	return registerCmd
}

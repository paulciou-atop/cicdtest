/*
Copyright © 2022 Atop NMS team

*/
package get

import (
	"encoding/json"
	"fmt"

	"nms/serviceswatcher/pkg/watcher"

	glog "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type GetOption struct {
	Name string //service's name
}

// NewCmdGet get specific service information
func NewCmdGet() *cobra.Command {

	var getCmd = &cobra.Command{
		Use:     "get",
		Short:   "get service's information",
		Example: "serviceswatcher get snmpscan",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			info, err := watcher.GetService(args[0])
			if err != nil {
				glog.Errorf("service %s did not find", args[0])
				return
			}
			jsonret, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				glog.Error("Convert result to JSON err: ", err)
				return
			}
			fmt.Println(string(jsonret))
		},
	}
	getCmd.SetUsageTemplate(" servicewatcher get [service-name]")

	return getCmd
}

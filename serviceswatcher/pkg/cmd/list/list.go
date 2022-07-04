/*
Copyright © 2022 Atop NMS team

*/
package list

import (
	"encoding/json"
	"fmt"

	"nms/serviceswatcher/pkg/watcher"

	"github.com/MakeNowJust/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCmdList List all running services
func NewCmdList() *cobra.Command {

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "list all running services",
		Long: heredoc.Doc(`
			List all running services. 
		`),
		Run: func(cmd *cobra.Command, args []string) {

			serviceList := watcher.ListServices()
			fmt.Println("serviceList ", serviceList)
			jsonret, err := json.MarshalIndent(serviceList, "", "  ")
			if err != nil {
				logrus.Error("Convert result to JSON err: ", err)
				return
			}
			fmt.Println(string(jsonret))

		},
	}

	return listCmd
}

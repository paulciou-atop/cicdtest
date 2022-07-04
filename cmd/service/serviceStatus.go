package cmd

import (
	"context"
	"encoding/json"
	"log"

	watcher "nms/api/v1/serviceswatcher"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StatusOptions struct {
	ServiceAddr string
	Name        string
}

func NewCmdStatus() *cobra.Command {
	o := &StatusOptions{}
	var statusCmd = &cobra.Command{
		Use:     "status",
		Short:   "Service commands: NMS service status",
		Example: "NMS service status snmpscan",
		Long: heredoc.Doc(` Atop NMS CLI
		This application is a tool to interface and control the Scanner service. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				//glog.Fatalf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			c := watcher.NewWatcherClient(conn)

			// Contact the server and print out its response.
			r, err := c.Get(context.Background(), &watcher.GetRequest{ServiceName: o.Name})
			if err != nil {
				//glog.Fatalf("could not connect snmp service: %v", err)
				return
			}

			jsonret, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				log.Println("Convert result to JSON err: ", err)
				return
			}

			cmd.Printf(string(jsonret))

			//cmd.Printf(r.Info.Name)
		},
	}

	statusCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")
	statusCmd.Flags().StringVarP(&o.Name, "name", "n", "", "Service name to get status")

	statusCmd.MarkFlagRequired("name")

	return statusCmd
}

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

type ListOptions struct {
	ServiceAddr string
}

func NewCmdList() *cobra.Command {
	o := &ListOptions{}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Service commands: NMS service list",
		Long: heredoc.Doc(` Atop NMS CLI
		This application is a tool to interface and control the Scanner service.. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			c := watcher.NewWatcherClient(conn)

			// Contact the server and print out its response.
			r, err := c.List(context.Background(), &watcher.Empty{})
			if err != nil {
				log.Fatalf("could not connect snmp service: %v", err)
				return
			}

			jsonret, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				log.Fatalf("could not transfer json format")
				return
			}
			cmd.Printf(string(jsonret))
		},
	}

	listCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")

	return listCmd
}

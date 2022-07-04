package cmd

import (
	"context"
	"encoding/json"

	scansrv "nms/api/v1/scanservice"

	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

type ScanStop struct {
	SessionId   string
	ServiceAddr string
	Output      string
}

func NewScanStopCmd() *cobra.Command {
	o := &ScanStop{}
	var scanStopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Scan commands: nmsctl scan stop",
		Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Println(err)
				return
			}
			defer conn.Close()
			c := scansrv.NewScanServiceClient(conn)

			// Contact the server and print out its response.
			r, err := c.StopScan(context.Background(), &scansrv.StopScanRequest{SessionId: o.SessionId})
			if err != nil {
				log.Println(err)
				return
			}

			switch {
			case o.Output == "json":
				jsonret, err := json.MarshalIndent(r, "", "  ")
				if err != nil {
					log.Fatalf("could not transfer json format")
					return
				}
				cmd.Printf(string(jsonret))
			default:
				yamlret, err := yaml.Marshal(r)
				if err != nil {
					log.Fatalf("could not transfer yaml format")
					return
				}
				cmd.Printf(string(yamlret))
			}
		},
	}

	scanStopCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")
	scanStopCmd.Flags().StringVarP(&o.SessionId, "sessionid", "s", "", "Service to get sessionID status")
	scanStopCmd.Flags().StringVarP(&o.Output, "output format", "o", "", "Output result with format, support json")

	scanStopCmd.MarkFlagRequired("sessionid")
	return scanStopCmd
}

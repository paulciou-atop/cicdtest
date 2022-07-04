package cmd

import (
	"context"
	"encoding/json"

	scansrv "nms/api/v1/scanservice"

	"log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

type ScanStatus struct {
	SessionId   string
	ServiceAddr string
	Output      string
}

func NewScanStatusCmd() *cobra.Command {
	o := &ScanStatus{}
	var scanStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Scan commands: nmsctl scan status",
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
			r, err := c.CheckStatus(context.Background(), &scansrv.CheckStatusRequest{SessionId: o.SessionId})
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
			case o.Output == "yaml":
				yamlret, err := yaml.Marshal(r)
				if err != nil {
					log.Fatalf("could not transfer yaml format")
					return
				}
				cmd.Printf(string(yamlret))
			default:
				t := table.NewWriter()
				t.SetStyle(table.StyleRounded)
				t.SetOutputMirror(cmd.OutOrStderr())
				t.AppendHeader(table.Row{"Status Scan Service"})
				var status = r.GetInfo()
				t.AppendRow(table.Row{status.GetStatus()})
				t.Render()
			}
		},
	}

	scanStatusCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")
	scanStatusCmd.Flags().StringVarP(&o.SessionId, "sessionid", "s", "", "Service to get sessionID status")
	scanStatusCmd.Flags().StringVarP(&o.Output, "output format", "o", "", "Output result with format, support json")

	scanStatusCmd.MarkFlagRequired("sessionid")
	return scanStatusCmd
}

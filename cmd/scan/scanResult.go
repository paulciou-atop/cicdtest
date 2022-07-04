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

type ScanResult struct {
	SessionId   string
	ServiceAddr string
	Page        int32
	Size        int32
	Output      string
}

func NewScanResultCmd() *cobra.Command {
	o := &ScanResult{}
	var scanResultCmd = &cobra.Command{
		Use:   "result",
		Short: "Scan commands: nmsctl scan result",
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
			r, err := c.GetResult(context.Background(), &scansrv.GetResultRequest{SessionId: o.SessionId, Page: o.Page, Size: o.Size})
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
				//status
				tw := table.NewWriter()
				tw.SetStyle(table.StyleRounded)
				tw.SetOutputMirror(cmd.OutOrStderr())
				var status = r.GetInfo()
				var statVal = status.GetMessage()
				if statVal == "no data for this session ID" {
					tw.AppendRow(table.Row{"Message", statVal})
				} else {
					tw.AppendRow(table.Row{"Message", "success"})
				}
				tw.Render()

				//content if no data
				t := table.NewWriter()
				t.SetStyle(table.StyleRounded)
				t.SetOutputMirror(cmd.OutOrStderr())
				t.AppendHeader(table.Row{"model", "mac address", "ip address", "netmask", "gateway", "hostname", "kernel", "ap", "firmware ver", "description", "device type", "scan time"})
				for _, c := range r.GetContent() {
					t.AppendRow(table.Row{c.GetModel(), c.GetMacAddress(), c.GetIpAddress(), c.GetNetmask(), c.GetGateway(), c.GetHostname(), c.GetKernel(), c.GetAp(), c.GetFirmwareVer(), c.GetDescription(), c.GetDeviceType(), c.GetScanTime()})
				}
				t.SetColumnConfigs([]table.ColumnConfig{
					{Name: "model", WidthMax: 10},
					{Name: "mac address", WidthMax: 11},
					{Name: "ip address", WidthMax: 10},
					{Name: "netmask", WidthMax: 10},
					{Name: "gateway", WidthMax: 8},
					{Name: "hostname", WidthMax: 8},
					{Name: "kernel", WidthMax: 6},
					{Name: "ap", WidthMax: 10},
					{Name: "firmware ver", WidthMax: 8},
					{Name: "description", WidthMax: 11},
					{Name: "device type", WidthMax: 6},
					{Name: "scan time", WidthMax: 11},
				})

				t.Render()
			}
		},
	}

	scanResultCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")
	scanResultCmd.Flags().StringVarP(&o.SessionId, "sessionid", "s", "", "Service to get sessionID status")
	scanResultCmd.Flags().Int32VarP(&o.Page, "page", "p", 1, "Result pagination")
	scanResultCmd.Flags().Int32VarP(&o.Size, "size", "z", 1, "Result size")
	scanResultCmd.Flags().StringVarP(&o.Output, "output format", "o", "", "Output result with format, support json")

	scanResultCmd.MarkFlagRequired("sessionid")
	return scanResultCmd
}

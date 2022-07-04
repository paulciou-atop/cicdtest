package cmd

import (
	"context"
	"encoding/json"

	scansrv "nms/api/v1/scanservice"

	"log"

	"github.com/gosnmp/gosnmp"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

type ScanOptions struct {
	ServiceAddr  string
	SnmpRange    string
	SnmpSettings SnmpOptions
	OIDsFile     string
	GwdServerIp  string
	Output       string
}

func NewScanStartCmd() *cobra.Command {
	o := &ScanOptions{}
	// scanCmd represents the scan command
	var scanStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Scan commands: nmsctl scan start",
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

			oids, _ := loadOIDsFile(o.OIDsFile)
			snmpVer := scansrv.SnmpSettings_SNMPVer(gosnmp.SnmpVersion(converSNMPVer(o.SnmpSettings.Version)))
			snmpSettings := scansrv.SnmpSettings{Port: o.SnmpSettings.Port, ReadCommunity: o.SnmpSettings.ReadCommunity, WriteCommunity: "", Version: snmpVer}

			// Contact the server and print out its response.
			r, err := c.StartScan(context.Background(), &scansrv.StartScanRequest{Range: o.SnmpRange, Oids: oids, SnmpSettings: &snmpSettings, ServerIp: o.GwdServerIp})
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
				t.AppendHeader(table.Row{"Session Id", "Message"})
				var info = r.GetInfo()
				t.AppendRow(table.Row{info.GetSessionId(), info.GetMessage()})
				t.SetColumnConfigs([]table.ColumnConfig{
					{Name: "Session Id", WidthMax: 15},
					{Name: "Message", WidthMax: 100},
					{Number: 5, WidthMax: 10},
				})
				t.Render()
			}
		},
	}

	scanStartCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")
	scanStartCmd.Flags().StringVarP(&o.SnmpRange, "snmp-range", "r", "172.17.0.0/24", "CIDR subnet defined")
	scanStartCmd.Flags().Int32VarP(&o.SnmpSettings.Port, "snmp-port", "p", 161, "Target SNMP agent port")
	scanStartCmd.Flags().StringVarP(&o.SnmpSettings.ReadCommunity, "snmp-community", "c", "public", "Target SNMP agent read community string")
	scanStartCmd.Flags().StringVarP(&o.SnmpSettings.Version, "snmp-version", "v", "v2", "SNMP version v1 | v2 | v3")
	scanStartCmd.Flags().StringVar(&o.OIDsFile, "snmp-oids-file", "", "OIDs plain text file")
	scanStartCmd.Flags().Int32Var(&o.SnmpSettings.Timeout, "snmp-timeout", 1, "SNMP timeout (second)")
	scanStartCmd.Flags().Int32Var(&o.SnmpSettings.Retries, "snmp-retries", 1, "SNMP retries time")
	scanStartCmd.Flags().StringVarP(&o.GwdServerIp, "gwd-server-ip", "s", "127.0.0.1", "GWD server ip")
	scanStartCmd.Flags().StringVarP(&o.Output, "output format", "o", "", "Output result with format, support json")

	viper.BindPFlag("snmp.port", Cmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("snmp.community", Cmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("snmp.version", Cmd.PersistentFlags().Lookup("snmp-version"))
	viper.BindPFlag("snmp.timeout", Cmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("snmp.retries", Cmd.PersistentFlags().Lookup("retries"))

	scanStartCmd.MarkFlagRequired("snmp-range")
	scanStartCmd.MarkFlagRequired("gwd-server-ip")

	return scanStartCmd
}

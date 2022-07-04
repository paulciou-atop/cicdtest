package cmd

import (
	"context"
	"encoding/json"
	"io"
	"log"

	scanV1 "nms/api/v1/snmpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/gosnmp/gosnmp"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

type GetOptions struct {
	ServiceAddr  string
	Target       string
	SnmpSettings SnmpOptions
	OIDsFile     string
	Output       string
}

// NewCmdGet create new get cobra command
func NewCmdGet() *cobra.Command {
	o := &GetOptions{}
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get OIDs' value of target",
		Long: heredoc.Doc(`
			Get OIDs from specfic target. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			defer conn.Close()
			c := scanV1.NewSnmpScanClient(conn)

			oids, _ := loadOIDsFile(o.OIDsFile)
			snmpVer := scanV1.SnmpSettings_SNMPVer(gosnmp.SnmpVersion(converSNMPVer(o.SnmpSettings.Version)))
			snmpSettings := scanV1.SnmpSettings{Port: o.SnmpSettings.Port, ReadCommunity: o.SnmpSettings.ReadCommunity, WriteCommunity: "", Version: snmpVer}

			// Contact the server and print out its response.
			r, err := c.Get(context.Background(), &scanV1.GetRequest{Target: o.Target, SnmpSettings: &snmpSettings, Oids: oids})
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStderr())
			t.SetStyle(table.StyleRounded)
			t.AppendHeader(table.Row{"Name", "Value", "Kind", "Oid"})
			def_format := false

			for {
				stream, err := r.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("error while reading stream: %v\n", err)
					return
				}

				switch {
				case o.Output == "json":
					jsonret, err := json.MarshalIndent(stream, "", "  ")
					if err != nil {
						log.Fatalf("could not transfer json format")
						return
					}
					cmd.Printf(string(jsonret))
				case o.Output == "yaml":
					// yaml encode not ready
					yamlret, err := yaml.Marshal(stream)
					if err != nil {
						log.Fatalf("could not transfer yaml format")
						return
					}
					cmd.Printf(string(yamlret))
				default:
					t.AppendRow(table.Row{stream.GetName(), stream.GetValue().GetStringValue(), stream.GetKind(), stream.GetOid()})
					def_format = true
				}
			}

			if def_format {
				t.Render()
			}
		},
	}

	getCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	getCmd.Flags().StringVarP(&o.Target, "target", "t", "127.0.0.1", "Target IP address")

	getCmd.Flags().Int32VarP(&o.SnmpSettings.Port, "port", "p", 161, "Target SNMP agent port")
	getCmd.Flags().StringVarP(&o.SnmpSettings.ReadCommunity, "community", "c", "public", "Target SNMP agent read community string")
	getCmd.Flags().StringVarP(&o.SnmpSettings.Version, "snmp-version", "v", "v2", "SNMP version v1 | v2 | v3")
	getCmd.Flags().Int32Var(&o.SnmpSettings.Timeout, "timeout", 1, "Snmp timeout (second)")
	getCmd.Flags().Int32Var(&o.SnmpSettings.Retries, "retries", 1, "SNMP retries time")
	getCmd.Flags().StringVar(&o.OIDsFile, "oids-file", "", "OIDs plain text file")
	getCmd.Flags().StringVarP(&o.Output, "output format", "o", "", "Output result with format, support json")

	viper.BindPFlag("serviceAddr", Cmd.PersistentFlags().Lookup("service-addr"))
	viper.BindPFlag("snmp.port", Cmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("snmp.community", Cmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("snmp.version", Cmd.PersistentFlags().Lookup("snmp-version"))
	viper.BindPFlag("snmp.timeout", Cmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("snmp.retries", Cmd.PersistentFlags().Lookup("retries"))

	getCmd.MarkFlagRequired("target")

	return getCmd
}

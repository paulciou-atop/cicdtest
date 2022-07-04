/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"nms/snmpscan/pkg/scan"

	"github.com/MakeNowJust/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ScanOptions struct {
	Range        string
	SnmpSettings SnmpOptions
	AtopDevices  bool
	Asynchronous bool
	OIDsFile     string
}

func NewScanCmd() *cobra.Command {
	o := &ScanOptions{}
	// scanCmd represents the scan command
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scan snmp agents",
		Long: heredoc.Doc(`
		Scan online devices which implement SNMP agent. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			oids, _ := loadOIDsFile(o.OIDsFile)

			opt := SnmpOptions2SnmpOption(&o.SnmpSettings)
			var results []scan.ScanResults
			var err error
			if o.Asynchronous {
				fmt.Println("execute async scan... ")
				id := scan.CreateSession()
				fmt.Println("session ID: ", id)
				scan.AsyncScan(id, o.Range, oids, opt, o.AtopDevices)
				fmt.Println("async scan done")
				return

			} else if o.AtopDevices {
				results, err = scan.ScanAtopDevices(o.Range, opt)
			} else {
				results, err = scan.ScanAgent(o.Range, oids, opt)
			}

			if err != nil {
				logrus.Error("Execute scan command fail : ", err)
				return
			}
			jsonret, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				logrus.Error("Convert result to JSON err: ", err)
				return
			}
			fmt.Println(string(jsonret))
		},
	}
	scanCmd.Flags().StringVarP(&o.Range, "range", "r", "172.17.0.0/24", "CIDR subnet defined")
	scanCmd.Flags().Int32VarP(&o.SnmpSettings.Port, "port", "p", 161, "Target SNMP agent port")
	scanCmd.Flags().StringVarP(&o.SnmpSettings.ReadCommunity, "community", "c", "public", "Target SNMP agent read community string")
	scanCmd.Flags().StringVarP(&o.SnmpSettings.Version, "snmp-version", "v", "v2", "SNMP version v1 | v2 | v3")
	scanCmd.Flags().StringVar(&o.OIDsFile, "oids-file", "", "OIDs plain text file")
	scanCmd.Flags().Int32Var(&o.SnmpSettings.Timeout, "timeout", 1, "Snmp timeout (second)")
	scanCmd.Flags().Int32Var(&o.SnmpSettings.Retries, "retries", 1, "SNMP retries time")
	scanCmd.Flags().BoolVar(&o.AtopDevices, "atop-devices", false, "Only scan atop devices")
	scanCmd.Flags().BoolVar(&o.Asynchronous, "async", false, "Asynchronous scan atop devices")

	viper.BindPFlag("snmp.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("snmp.community", rootCmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("snmp.version", rootCmd.PersistentFlags().Lookup("snmp-version"))
	viper.BindPFlag("snmp.timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("snmp.retries", rootCmd.PersistentFlags().Lookup("retries"))
	scanCmd.MarkFlagRequired("range")
	return scanCmd
}

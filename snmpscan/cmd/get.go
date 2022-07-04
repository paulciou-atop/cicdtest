/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"nms/snmpscan/pkg/scan"
	"nms/snmpscan/pkg/snmp"

	"github.com/MakeNowJust/heredoc"
	glog "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GetOptions struct {
	Target       string
	SnmpSettings SnmpOptions
	OIDsFile     string
}

// NewCmdGet create new get cobra command
func NewCmdGet() *cobra.Command {
	o := &GetOptions{}
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information (snmp|session) ",
		Long: heredoc.Doc(`
			Get some kind of information, please refer --help for details. 
		`),
	}

	var snmpCmd = &cobra.Command{
		Use:     "snmp [target]",
		Short:   "get target's snmp infmormation",
		Example: "snmpscan get 192.168.1.11",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			oids, _ := loadOIDsFile(o.OIDsFile)

			opt := SnmpOptions2SnmpOption(&o.SnmpSettings)

			results, err := snmp.Get(args[0], oids, opt)
			if err != nil {
				glog.Error("Execute get command fail : ", err)
				return
			}
			var pdus []snmp.ResultItem
			for _, v := range results {

				pdus = append(pdus, v)
			}
			if err != nil {
				glog.Error("Convert result to JSON err: ", err)
				return
			}
			jsonret, err := json.MarshalIndent(pdus, "", "  ")
			if err != nil {
				glog.Error("Convert result to JSON err: ", err)
				return
			}
			fmt.Println(string(jsonret))
		},
	}

	var sessionCmd = &cobra.Command{
		Use:     "session [session-id]",
		Short:   "get session's result with specifc session-id",
		Example: "snmpscan get session ss:f9083d4d-5d18-492b-ac6d-1a7506fa3ab0",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ret, err := scan.GetAsyncResult(args[0])
			if err != nil {
				glog.Errorf("Execute get session %s fail : %+v", args[0], err)
				return
			}
			jsonret, err := json.MarshalIndent(ret, "", "  ")
			if err != nil {
				glog.Error("Convert result to JSON err: ", err)
				return
			}
			fmt.Println(string(jsonret))
		},
	}

	getCmd.AddCommand(snmpCmd)
	getCmd.AddCommand(sessionCmd)

	snmpCmd.Flags().Int32VarP(&o.SnmpSettings.Port, "port", "p", 161, "Target SNMP agent port")
	snmpCmd.Flags().StringVarP(&o.SnmpSettings.ReadCommunity, "community", "c", "public", "Target SNMP agent read community string")
	snmpCmd.Flags().StringVarP(&o.SnmpSettings.Version, "snmp-version", "v", "v2", "SNMP version v1 | v2 | v3")
	snmpCmd.Flags().Int32Var(&o.SnmpSettings.Timeout, "timeout", 1, "Snmp timeout (second)")
	snmpCmd.Flags().Int32Var(&o.SnmpSettings.Retries, "retries", 1, "SNMP retries time")
	snmpCmd.Flags().StringVar(&o.OIDsFile, "oids-file", "", "OIDs plain text file")

	viper.BindPFlag("snmp.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("snmp.community", rootCmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("snmp.version", rootCmd.PersistentFlags().Lookup("snmp-version"))
	viper.BindPFlag("snmp.timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("snmp.retries", rootCmd.PersistentFlags().Lookup("retries"))

	return getCmd
}

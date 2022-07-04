/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"nms/snmpscan/pkg/snmp"

	"github.com/MakeNowJust/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type WalkOption struct {
	RootOid      string
	Target       string
	SnmpSettings SnmpOptions
}

func NewCmdWalk() *cobra.Command {
	o := &WalkOption{}
	//wakl command
	var walkCmd = &cobra.Command{
		Use:   "walk",
		Short: "get subtree",
		Long: heredoc.Doc(`
		Scan online devices which implement SNMP agent. 
		`),
		Run: func(cmd *cobra.Command, args []string) {

			opt := SnmpOptions2SnmpOption(&o.SnmpSettings)
			var results []snmp.ResultItem
			var err error
			results, err = snmp.Walk(o.Target, o.RootOid, opt)

			if err != nil {
				logrus.Error("Execute get command fail : ", err)
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
	walkCmd.Flags().StringVarP(&o.RootOid, "root-oid", "r", ".1.3.6.1.2.1.2.2.1.6", "Root oid")
	walkCmd.Flags().StringVarP(&o.Target, "target", "t", "localhost", "Target IP Address")
	walkCmd.Flags().Int32VarP(&o.SnmpSettings.Port, "port", "p", 161, "Target SNMP agent port")
	walkCmd.Flags().StringVarP(&o.SnmpSettings.ReadCommunity, "community", "c", "public", "Target SNMP agent read community string")
	walkCmd.Flags().StringVarP(&o.SnmpSettings.Version, "snmp-version", "v", "v2", "SNMP version v1 | v2 | v3")
	walkCmd.Flags().Int32Var(&o.SnmpSettings.Timeout, "timeout", 1, "Snmp timeout (second)")
	walkCmd.Flags().Int32Var(&o.SnmpSettings.Retries, "retries", 1, "SNMP retries time")

	viper.BindPFlag("snmp.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("snmp.community", rootCmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("snmp.version", rootCmd.PersistentFlags().Lookup("snmp-version"))
	viper.BindPFlag("snmp.timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("snmp.retries", rootCmd.PersistentFlags().Lookup("retries"))
	walkCmd.MarkFlagRequired("target")
	return walkCmd
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"nms/snmpscan/pkg/snmp"

	"github.com/MakeNowJust/heredoc"
	glog "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DescribeOption struct {
	Target       string
	SnmpSettings SnmpOptions
}

func NewCmdDescribe() *cobra.Command {
	o := &WalkOption{}
	//wakl command
	var descrCmd = &cobra.Command{
		Use:   "describe [target]",
		Short: "describe target's information",
		Long: heredoc.Doc(`
		get details information snmp information from the target. 
		`),
		Example: "snmpscan describe 192.169.1.2",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			opt := SnmpOptions2SnmpOption(&o.SnmpSettings)

			result, err := snmp.Describe(args[0], opt)
			if err != nil {
				glog.Error("Execute describe command fail : ", err)
				return
			}
			jsonret, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				glog.Error("Convert result to JSON err: ", err)
				return
			}
			fmt.Println(string(jsonret))

		},
	}

	descrCmd.PersistentFlags().Int32VarP(&o.SnmpSettings.Port, "port", "p", 161, "Target SNMP agent port")
	descrCmd.PersistentFlags().StringVarP(&o.SnmpSettings.ReadCommunity, "community", "c", "public", "Target SNMP agent read community string")
	descrCmd.PersistentFlags().StringVarP(&o.SnmpSettings.Version, "snmp-version", "v", "v2", "SNMP version v1 | v2 | v3")
	descrCmd.PersistentFlags().Int32Var(&o.SnmpSettings.Timeout, "timeout", 1, "Snmp timeout (second)")
	descrCmd.PersistentFlags().Int32Var(&o.SnmpSettings.Retries, "retries", 1, "SNMP retries time")

	viper.BindPFlag("snmp.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("snmp.community", rootCmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("snmp.version", rootCmd.PersistentFlags().Lookup("snmp-version"))
	viper.BindPFlag("snmp.timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("snmp.retries", rootCmd.PersistentFlags().Lookup("retries"))
	descrCmd.MarkFlagRequired("target")
	return descrCmd
}

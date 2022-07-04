package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(NewCmdGet())
	Cmd.AddCommand(NewScanCmd())
	Cmd.AddCommand(NewCmdWalk())
}

var Cmd = &cobra.Command{
	Use:   "snmpscan",
	Short: "SNMP commands: nmsctl snmpscan --help",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
	Run: runHelp,
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(NewCmdGetServerIp())
	Cmd.AddCommand(NewCmdBeep())
	Cmd.AddCommand(NewCmdScan())
	Cmd.AddCommand(NewCmdStop())
	Cmd.AddCommand(NewCmdSessionData())
}

var Cmd = &cobra.Command{
	Use:   "udp",
	Short: "UDP commands: nmsctl udp --help",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
}

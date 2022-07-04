package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(NewCmdList())
	Cmd.AddCommand(NewCmdStatus())
	Cmd.AddCommand(NewCmdRegister())
}

var Cmd = &cobra.Command{
	Use:   "service",
	Short: "Service commands: nmsctl service --help",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
}

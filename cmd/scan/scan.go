package cmd

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan commands: nmsctl scan --help",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
}

func init() {
	Cmd.AddCommand(NewScanStartCmd())
	Cmd.AddCommand(NewScanStatusCmd())
	Cmd.AddCommand(NewScanStopCmd())
	Cmd.AddCommand(NewScanResultCmd())
}

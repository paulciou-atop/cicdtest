/*
Package root implements root CLI command

*/
package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	configCmd "nms/config/pkg/cmd/config"
	getCmd "nms/config/pkg/cmd/get"
	runCmd "nms/config/pkg/cmd/run"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command> <subcommand> [flags]",
		Short: "NMS config CLI",
		Long:  `Device configuration tools with NMS from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ config config 00-00-00-12-12-34 --ip=110.34.14.57
		`),
		Run: runHelp,
	}
	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetFlagErrorFunc(rootFlagErrorHandler)

	// Sub-commands
	cmd.AddCommand(runCmd.NewCmdRun())
	cmd.AddCommand(configCmd.NewCmdConfig())
	cmd.AddCommand(getCmd.NewCmdGet())
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

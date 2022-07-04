/*
Package root implements root CLI command

*/
package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	runCmd "nms/testing/dummyconfiger/pkg/cmd/run"
)

func rootFlagErrorHandler(cmd *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}
	return err
}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nconfig <command> <subcommand> [flags]",
		Short: "NMS config CLI",
		Long:  `Device configuration tools with NMS from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ nconfig config 00-00-00-12-12-34 --ip=110.34.14.57
		`),
		Run: runHelp,
	}
	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetFlagErrorFunc(rootFlagErrorHandler)

	// Sub-commands
	cmd.AddCommand(runCmd.NewCmdRun())
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

/*
Package root implements root CLI command

*/
package root

import (
	listCmd "nms/inventory/pkg/cmd/list"
	runCmd "nms/inventory/pkg/cmd/run"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inventory <command> <subcommand> [flags]",
		Short: "Manage inventories",
		Long:  `Work with inventories.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ inventory list
		`),
		Run: runHelp,
	}
	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.AddCommand(runCmd.NewCmdRun())
	cmd.AddCommand(listCmd.NewCmdList())
	// Sub-commands

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

/*
Package cmd implements bunch of functions about CLI

*/
package cmd

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "messaging",
	Short: "messaging publish or subscribe messages",
	Long: heredoc.Doc(`
		messaging publish or subscribe messages.

	`),
	Run: runHelp,
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Default CLI execute err: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewCmdPub())
	rootCmd.AddCommand(NewCmdSub())
}

// NewDefaultCommand creates the `snmp` command with default arguments
func NewDefaultCommand() *cobra.Command {

	return nil
}

// NewDefaultCommandWithArgs
//func NewDefaultCommandWithArgs(o Options)

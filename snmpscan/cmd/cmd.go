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
	Use:   "snmpscan",
	Short: "snmpscan scan snmp devices/agents",
	Long: heredoc.Doc(`
		snmpscan scan snmp devices/agents.

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
	rootCmd.AddCommand(NewCmdGet())
	rootCmd.AddCommand(NewScanCmd())
	rootCmd.AddCommand(NewRunCmd())
	rootCmd.AddCommand(NewCmdWalk())
	rootCmd.AddCommand(NewCmdDescribe())
	rootCmd.AddCommand(NewCmdTest())

}

// NewDefaultCommand creates the `snmp` command with default arguments
func NewDefaultCommand() *cobra.Command {

	return nil
}

// NewDefaultCommandWithArgs
//func NewDefaultCommandWithArgs(o Options)

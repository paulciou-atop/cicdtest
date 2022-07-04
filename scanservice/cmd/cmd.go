/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Run scanservice service",
	Long:  heredoc.Doc("Running gRPC & HTTP/HTTPS server"),
	Run:   runHelp,
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		glog.Error("Default CLI execute err: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(Run())
}

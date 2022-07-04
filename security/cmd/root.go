package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "security",
}

func init() {
	RootCmd.AddCommand(NewRunCmd())
	RootCmd.AddCommand(NewRootCrtCmd())
	RootCmd.AddCommand(NewSrvCrtCmd())
	RootCmd.AddCommand(NewSrvKeyCmd())
}
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

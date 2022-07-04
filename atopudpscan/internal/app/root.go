package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//var address = "localhost:" + config.GetgrpcPort()
const account = "admin"
const password = "default"

var RootCmd = &cobra.Command{
	Use: "atopudpscan",
}

func init() {
	RootCmd.AddCommand(NewRunCmd())
	RootCmd.AddCommand(GwdCmd)
	RootCmd.AddCommand(deviceCmd)
	RootCmd.AddCommand(ConfigCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

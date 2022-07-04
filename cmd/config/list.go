package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands: nmsctl config list",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func list() {
	v1 := viper.New()
	v1.SetConfigFile("config.yaml")
	v1.ReadInConfig()
	fmt.Println("username:", v1.GetString("username"))
	fmt.Println("password:", v1.GetString("password"))
}

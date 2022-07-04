package main

import (
	"fmt"
	"os"

	config "nms/cmd/config"
	scan "nms/cmd/scan"
	service "nms/cmd/service"
	snmpscan "nms/cmd/snmpscan"
	udpscan "nms/cmd/udpscan"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "nmsctl",
	Version: "v1.0.0-alpha.1",
	Short:   "Atop NMS CLI commands: NMS",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
}

func main() {
	addCommands()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.NMS.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".NMS")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func addCommands() {
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(scan.Cmd)
	rootCmd.AddCommand(service.Cmd)
	rootCmd.AddCommand(snmpscan.Cmd)
	rootCmd.AddCommand(udpscan.Cmd)
}

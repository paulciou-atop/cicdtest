package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func init() {
	Cmd.AddCommand(listCmd)
}

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Config commands: nmsctl config",
	Long: `Atop NMS CLI
This application is a tool to interface and control the Scanner service.`,
	Run: func(cmd *cobra.Command, args []string) {
		writeTemplate()
	},
}

type TemplateJSONType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func writeTemplate() {
	var templateJSON TemplateJSONType

	if viper.GetString("username") != "" {
		templateJSON.Username = viper.GetString("username")
	}
	if viper.GetString("password") != "" {
		templateJSON.Password = viper.GetString("password")
	}

	data, err := yaml.Marshal(templateJSON)

	if err != nil {

		log.Fatal(err)
	}

	err2 := ioutil.WriteFile("config.yaml", data, 0)

	if err2 != nil {

		log.Fatal(err2)
	}

	fmt.Println("data written")
}

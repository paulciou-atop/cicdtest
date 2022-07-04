/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"nms/snmpscan/pkg/store"

	"github.com/spf13/cobra"
)

func NewCmdTest() *cobra.Command {

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "testing stuff",

		Run: func(cmd *cobra.Command, args []string) {
			err := store.ExampleRedisClient()
			if err != nil {
				fmt.Println("Testing redis error: ", err)
				return
			}
			fmt.Println("Testing done")
		},
	}

	return testCmd
}

package main

import (
	"nms/inventory/pkg/cmd/root"

	"github.com/sirupsen/logrus"
)

func main() {

	rootCommand := root.NewCmdRoot()

	if err := rootCommand.Execute(); err != nil {
		logrus.Error("Command fail: ", err)
		return
	}
}

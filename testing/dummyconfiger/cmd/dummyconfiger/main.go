package main

import (
	"nms/testing/dummyconfiger/pkg/cmd/root"

	log "github.com/sirupsen/logrus"
)

func main() {

	rootCommand := root.NewCmdRoot()

	if err := rootCommand.Execute(); err != nil {
		log.Error("Command fail: ", err)
		return
	}
}

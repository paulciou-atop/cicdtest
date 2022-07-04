package main

import (
	"nms/snmpscan/cmd"
	_ "nms/snmpscan/pkg/config"
)

func main() {
	//server.RunGRPCServer()
	cmd.Execute()
}

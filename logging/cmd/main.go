package main

import (
	AtopSyslog "atopsyslog/syslog"
	"log"
)

func main() {
	server := AtopSyslog.NewServer("syslog.log")
	err := server.Run("0.0.0.0:514")
	if err != nil {
		log.Fatal(err)
	}
}

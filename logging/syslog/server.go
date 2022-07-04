package AtopSyslog

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/mcuadros/go-syslog.v2"
)

type Server struct {
	filename string
}

func NewServer(filename string) *Server {
	return &Server{filename}
}

func (s *Server) Run(network string) error {

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	fmt.Println("Syslog Server start..")
	err := server.ListenUDP(network)
	if err != nil {
		return err
	}
	err = server.Boot()
	if err != nil {
		return err
	}
	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
			writeFile(s.filename, logParts)
		}

	}(channel)
	fmt.Println("Syslog Server Listening..")
	server.Wait()
	return nil
}

func writeFile(filename string, value interface{}) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	js, _ := json.Marshal(value)
	f.Write(js)
	f.Write([]byte("\n"))
	f.Close()
}

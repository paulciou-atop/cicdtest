package gwd_test

import (
	"log"

	"nms/atopudpscan/pkg/gwd/server"
	"os"
	"os/signal"
	"testing"
)

func TestServer(t *testing.T) {
	s := server.NewAtopUdpServer("0.0.0.0")
	err := s.Run()
	if err != nil {
		t.Error(err)
		return
	}
	go func() {
		for {
			v, err := s.GetReceiveData()
			if err == nil {
				log.Println(string(v))

			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.Stop()
}

func TestSaveToDatabase(t *testing.T) {
	s := server.NewAtopUdpServer("0.0.0.0")
	go s.Run()
	go func() {
		for {
			v, err := s.GetReceiveData()
			if err == nil {
				log.Println(string(v))
				err := s.SaveToDatabase(v)
				if err != nil {
					t.Error(err)
					return
				}
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.Stop()
}
func TestGetDataFromDatabase(t *testing.T) {
	s := server.NewAtopUdpServer("0.0.0.0")
	b, err := s.GetDataFromDatabase()
	if err == nil {
		log.Println(string(b))
	} else {
		t.Error(err)
	}
}

func TestCleanDatabase(t *testing.T) {
	s := server.NewAtopUdpServer("0.0.0.0")
	err := s.CleanDatabase()
	if err != nil {
		t.Error(err)
	}
}

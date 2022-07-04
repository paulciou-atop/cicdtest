package AtopSyslog

import (
	"testing"
)

func TestServer(t *testing.T) {

	s := NewServer("test.log")
	err := s.Run("192.168.4.21:514")
	if err != nil {
		t.Error(err.Error())
	}

}

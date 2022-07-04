package AtopSyslog

import (
	"testing"
)

func TestSendingMessage(t *testing.T) {

	sysLog, err := DialLogger("udp", "192.168.4.21:514", LOG_DEBUG|LOG_SYSLOG, "demotag")
	if err != nil {
		t.Error(err)
	}
	//	fmt.Fprintf(sysLog, "This is a daemon warning with demotag.")
	sysLog.Write([]byte("test message"))
	sysLog.Close()
}

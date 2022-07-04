package session

import (
	"fmt"
	pg "nms/lib/pgutils"

	"github.com/sirupsen/logrus"
)

// InitDatabaseTables use this function to init pgutils.IClient's table schema
func InitDatabaseTables(client pg.IClient) error {
	err := client.CreateTable(&ConfigTask{}, pg.CreateTableOpt{IfNotExists: true})
	if err != nil {
		logrus.Error("create config task table fail: ", err)

	}
	err = client.CreateTable(&ConfigSession{}, pg.CreateTableOpt{IfNotExists: true})
	if err != nil {
		logrus.Error("create session state table fail: ", err)

	}

	return nil
}

// SessionTable
type ConfigSession struct {
	SessionID string `pg:",pk"` //primary key
	State     string
	StartedAt string
	EndedAt   string
	Message   string
}

func (s ConfigSession) String() string {
	return fmt.Sprintf("Session-ID:%s\nState:%s\ntime:%s - %s\n",
		s.SessionID,
		s.State,
		s.StartedAt,
		s.EndedAt)
}

// ConfigTask, config task is
type ConfigTask struct {
	Id           int64
	SessionID    string
	DeviceID     string
	DevicePath   string
	ConfigHash   string
	FaildOptions []string
}

func (t ConfigTask) IsSuccess() bool {
	if len(t.FaildOptions) <= 0 {
		return true
	}
	return false
}

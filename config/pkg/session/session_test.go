package session_test

import (
	"context"
	"nms/api/v1/common"
	"nms/api/v1/configer"
	"nms/api/v1/devconfig"
	"nms/config/internal/repo"
	"nms/config/pkg/session"
	"testing"
	"time"
)

var fakeConfigPayload = map[string]interface{}{
	"ip":       "192.168.13.1",
	"mask":     "255.255.255.0",
	"gateway":  "192.168.13.254",
	"account":  "user",
	"password": "password",
}

func TestSession(t *testing.T) {
	r, err := repo.GetRepo(context.Background())
	if err != nil {
		t.Error(err)
	}

	db := r.DB()

	defer r.CleanUp()

	pgdb, err := db.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS config_sessions`)
	if err != nil {
		t.Error(err)
	}

	_, err = pgdb.Exec(`DROP TABLE IF EXISTS config_tasks`)
	if err != nil {
		t.Error(err)
	}

	session.InitDatabaseTables(db)
	s := session.NewSession(context.TODO())
	// Request coming
	// save config
	// configMetric1 := config.NewConfigMetric("network", fakeConfigPayload)
	// got new testSession
	testSession := <-s
	// got config
	response := configer.ConfigerResponse{
		Session: &devconfig.SessionState{
			Id:          testSession.SessionID(),
			StartedTime: testSession.State.StartedAt,
			State:       "success",
			EndedTime:   time.Now().Format(session.TIME_LAYOUT),
		},
		Device: &common.DeviceIdentify{
			DeviceId:   "test_dev_id",
			DevicePath: "test_dev_path",
		},
		ConfigResults: []*devconfig.ConfigResult{
			{
				Protocol:   "dummyconfiger",
				Kind:       "testing",
				Hash:       "hash123",
				FailFields: []string{},
			},
		},
	}

	tasks, err := session.MarshalConfigTasks(&response)
	if err != nil {
		t.Errorf("marshal response to task fail: %v", err)
		return
	}

	testSession.AddTasks(tasks)
	testSession.Done()
	testSession.WriteToStore(db)

	report, err := session.GetConfigSession(db, testSession.SessionID())
	if err != nil {
		t.Error(err)
	}
	// Dump session report
	t.Logf(report.String())
	// testing case
	if report.State != "success" {
		t.Errorf("should success but not %s", report.State)
	}
	storeTasks, err := session.GetConfigTasks(db, testSession.SessionID())
	if err != nil {
		t.Error(err)
	}
	if len(storeTasks) != 1 {
		t.Errorf("should has one task but %d", len(storeTasks))
	}
}

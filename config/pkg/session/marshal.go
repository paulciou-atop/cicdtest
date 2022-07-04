/*
 All marshal functions
*/
package session

import (
	"fmt"
	configerAPI "nms/api/v1/configer"
	"nms/api/v1/devconfig"
)

// UnMarshalSessionState   session.ConfigSession -> devconfig.SessionState
func UnMarshalSessionState(cs *ConfigSession) *devconfig.SessionState {
	return &devconfig.SessionState{
		Id:          cs.SessionID,
		State:       cs.State,
		StartedTime: cs.StartedAt,
		EndedTime:   cs.EndedAt,
		Message:     cs.Message,
		// TODO Topic: s.State.,
	}
}

// MarshalConfigSession  configer.ConfigerResponse -> ConfigSession
func MarshalConfigSession(res *configerAPI.ConfigerResponse) (ConfigSession, error) {
	empytConfigSession := new(ConfigSession)
	if res.Session == nil {
		return *empytConfigSession, fmt.Errorf("marshal fail : got nil session")
	}

	if res.Device == nil {
		return *empytConfigSession, fmt.Errorf("marshal fail : got nil device")
	}
	return ConfigSession{
		SessionID: res.Session.Id,
		State:     res.Session.State,
		StartedAt: res.Session.StartedTime,
		EndedAt:   res.Session.EndedTime,
		Message:   res.Session.Message,
	}, nil
}

// MarshalConfigTasks configerAPI.ConfigerResponse -> []ConfigTask
func MarshalConfigTasks(res *configerAPI.ConfigerResponse) ([]ConfigTask, error) {
	tasks := []ConfigTask{}
	if res.Session == nil {
		return []ConfigTask{}, fmt.Errorf("marshal fail : got nil session")
	}
	if res.Device == nil {
		return []ConfigTask{}, fmt.Errorf("marshal fail : got nil device")
	}
	for _, r := range res.ConfigResults {
		t := ConfigTask{
			SessionID:    res.Session.Id,
			DeviceID:     res.Device.DeviceId,
			DevicePath:   res.Device.DevicePath,
			ConfigHash:   r.Hash,
			FaildOptions: r.FailFields,
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

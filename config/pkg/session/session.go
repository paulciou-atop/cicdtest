package session

import (
	"context"
	"encoding/json"
	"fmt"

	"nms/lib/pgutils"
	pg "nms/lib/pgutils"
	mq "nms/messaging"
	"sync"
	"time"

	"github.com/google/uuid"
)

const TIME_LAYOUT = pgutils.TIME_LAYOUT

func timeString(t time.Time) string {
	return t.Local().Format(TIME_LAYOUT)
}

func NewConfigSession() ConfigSession {
	configSession := ConfigSession{
		SessionID: uuid.New().String(),
		StartedAt: time.Now().Format(TIME_LAYOUT),
	}
	configSession.Running()
	return configSession
}

func (s *ConfigSession) Running() {
	s.State = "running"
}

func (s *ConfigSession) Done(success bool) {
	if success {
		s.State = "success"
	} else {
		s.State = "fail"
	}
	s.EndedAt = time.Now().Format(TIME_LAYOUT)
}

func (s *ConfigSession) Cancel() {
	s.State = "cancel"
	s.EndedAt = time.Now().Format(TIME_LAYOUT)
}

// Session Session data module defined
/* add task :
	session.AddTask <- &ConfigTask{}
 done:
  // value dosen't matter
  session.Finished <- 1
*/
type Session struct {
	State    ConfigSession
	Tasks    []ConfigTask
	addTasks chan<- []ConfigTask
	update   chan<- ConfigSession
	finished chan<- time.Time
}

//SessionID return session id
func (s Session) SessionID() string {
	return s.State.SessionID
}

var taskMux sync.Mutex

// AddTasks
func (s *Session) AddTasks(tasks []ConfigTask) {
	var newTasks []ConfigTask
	// check duplicate tasks
	// if task have same hash and deviceID, it's duplicate
	if len(tasks) <= 0 {
		return
	}

	taskMux.Lock()
	for _, t := range tasks {
		key := t.DeviceID + t.ConfigHash
		exist := false
		for _, task := range s.Tasks {
			if key == task.DeviceID+task.ConfigHash {
				exist = true
			}
		}
		if !exist {
			newTasks = append(newTasks, t)
		}
	}
	taskMux.Unlock()

	s.addTasks <- newTasks
}

func (s *Session) Update(cs ConfigSession) {
	s.update <- cs
}

func (s *Session) Done() {
	s.finished <- time.Now()
	close(s.update)
	close(s.addTasks)
	close(s.finished)
}

func (s *Session) WriteToStore(store pg.IDataAccess) error {
	// check state exist?
	var sessions []ConfigSession
	store.Query(&sessions, pg.QueryExpr{
		Expr: "session_id = ?", Value: s.SessionID(),
	})
	if len(sessions) > 0 {
		if err := store.Update(&s.State); err != nil {
			return err
		}
	} else {
		if err := store.Insert(&s.State); err != nil {
			return err
		}
	}
	if len(s.Tasks) > 0 {
		if err := store.Insert(&s.Tasks); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) Publish(m mq.IMQClient) error {
	jsonret, err := json.MarshalIndent(s.State, "", "  ")
	if err != nil {
		return err
	}
	return m.Publish("config.config", string(jsonret))
}

// return first 6 character of hash
func shortHash(hash string) string {
	i := 0
	for j := range hash {
		if i == 6 {
			return hash[:j]
		}
		i++
	}
	return hash
}

// DumpConfigTasks prety string of config tasks
func DumpConfigTasks(tasks []ConfigTask) string {
	var result = fmt.Sprintf("Total tasks: %d\n", len(tasks))
	for _, v := range tasks {
		result += fmt.Sprintf("Device:%s %s\nConfig:%s\n", v.DeviceID, v.DevicePath, shortHash(v.ConfigHash))
		result += fmt.Sprintf("fail options: %v", v.FaildOptions)
	}
	return result

}

// NewSession create a new session struct
// example add task
// session := NewSession(ctx,db)
// session.AddTask <- task
// session.Finished <- time.Now()
func NewSession(ctx context.Context) chan *Session {

	// output session chan
	c := make(chan *Session)
	// add task channel
	go func() {
		addtask := make(chan []ConfigTask)
		finished := make(chan time.Time)
		update := make(chan ConfigSession)
		session := &Session{
			State:    NewConfigSession(),
			Tasks:    []ConfigTask{},
			finished: finished,
			addTasks: addtask,
			update:   update,
		}

		c <- session

		var success = true
		// session main routine
		for {
			select {
			case newTasks := <-addtask:
				// add task
				if len(newTasks) < 1 || newTasks == nil {
					continue
				}
				session.Tasks = append(session.Tasks, newTasks...)
				for _, t := range newTasks {
					if !t.IsSuccess() {
						success = false
					}
				}

			case sc := <-update:
				session.State = sc

			case <-finished:
				// do finished
				session.State.Done(success)

				return
			case <-ctx.Done():
				// do cancel
				session.State.Cancel()
				return
			}
		}
	}()

	return c
}

func GetAllConfigSession(ctx context.Context, querier pg.Querier, state ...string) ([]ConfigSession, error) {

	var sessions []ConfigSession

	if err := querier.Query(&sessions); err != nil {
		return []ConfigSession{}, err
	}

	return sessions, nil
}

// GetSessionReport get session report
func GetConfigSession(db pg.Querier, id string) (ConfigSession, error) {
	var sessions []ConfigSession
	if err := db.Query(&sessions, pg.QueryExpr{
		Expr:  "session_id = ?",
		Value: id,
	}); err != nil {
		return ConfigSession{}, err
	}
	if len(sessions) != 1 {
		return ConfigSession{}, fmt.Errorf("multiple session state with same id %s", id)
	}
	return sessions[0], nil
}

func GetConfigTasks(db pg.Querier, sessionID string) ([]ConfigTask, error) {
	var tasks []ConfigTask
	if err := db.Query(&tasks, pg.QueryExpr{
		Expr:  "session_id = ?",
		Value: sessionID,
	}); err != nil {
		return []ConfigTask{}, err
	}

	return tasks, nil
}

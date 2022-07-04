package store

import (
	"fmt"
	"strings"
	"time"

	pg "github.com/Atop-NMS-team/pgutils"
	log "github.com/sirupsen/logrus"
)

type scansession struct {
	client pg.IDBClient
	id     string
}

func NewSession(sessionID string) (*scansession, error) {
	//prepare store
	c, err := pg.NewClient()

	if err != nil {
		log.Error("create postgres client fail: ", err)
		return nil, err

	}
	ss := &scansession{
		client: c,
		id:     sessionID,
	}

	return ss, nil
}

func (ss *scansession) Close() {
	if ss.client != nil {
		ss.client.Close()
	}
}

type expr = pg.QueryExpr

// updateSessionState in the db, session state like this gwe:running|snmp:running
// this handily function help to create correct format of seesion state
func updateSessionState(ori string, newstate string) string {
	data := map[string]string{}
	states := strings.Split(ori, "|")

	for _, i := range states {
		kv := strings.Split(i, ":")
		if len(kv) < 2 {
			data["gwd"] = ""
		} else {
			data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	data["snmp"] = newstate
	return fmt.Sprintf("gwd:%s|snmp:%s", data["gwd"], data["snmp"])
}

// StoreSession store session to postgres
func (ss *scansession) StoreSession(data *pg.DeviceSession) error {
	c := ss.client
	sessionid := ss.id
	// make sure table exist
	c.CreateTable(&pg.DeviceSession{}, pg.CreateTableOpt{IfNotExists: true})

	var findSessions []pg.DeviceSession
	err := c.Query(&findSessions, pg.QueryExpr{Expr: "session_id = ?", Value: sessionid})
	if err != nil {
		log.Errorf("query session id fail: %v", err)
		return err
	}

	if len(findSessions) < 1 {
		// error session not found
		log.Errorf("query session id not found %s", sessionid)
		return fmt.Errorf("session id %s not found", sessionid)
	}

	newdata := findSessions[0]
	newdata.State = updateSessionState(newdata.State, data.State)
	newdata.LastUpdatedTime = time.Now().String()
	err = c.Update(&newdata)
	if err != nil {
		log.Errorf("update session id fail: %v", err)
		return err
	}
	return nil
}

func (ss *scansession) StoreDeviceResults(scanDevs []pg.DeviceResult) error {

	c := ss.client
	sessionID := ss.id

	c.CreateTable(&pg.DeviceResult{}, pg.CreateTableOpt{IfNotExists: true})

	for _, dev := range scanDevs {
		// step1 find record which match sessionID and MAC if not create a new one
		var queryResults []pg.DeviceResult
		err := c.Query(&queryResults,
			expr{Expr: "session_id = ?", Value: sessionID},
			expr{Expr: "mac_address = ?", Value: dev.MacAddress},
		)
		if err != nil {
			log.Error("querry fail:", err)
			return err
		}
		if len(queryResults) < 1 {
			// not found! insert one
			err := c.Insert(&dev)
			if err != nil {
				log.Error("insert device result fail: ", err)
				return err
			}
		}
		if len(queryResults) == 1 {
			// have dev informaiton update it
			newDev := queryResults[0]
			newDev.Description = dev.Description
			newDev.FirmwareVer = dev.FirmwareVer
			newDev.MacAddress = dev.MacAddress
			newDev.Kernel = dev.Kernel
			newDev.Model = dev.Model
			newDev.IpAddress = dev.IpAddress
			err = c.Update(newDev)
			if err != nil {
				log.Errorf("update device mac:%s result fail: %v", dev.MacAddress, err)
				return err
			}
		}
	}

	return nil
}

package service

import (
	"log"
	"time"

	"github.com/Atop-NMS-team/pgutils"

	"nms/scanservice/api/utils"
)

func ConnectDB() map[string]any {
	// create new client
	c, err1 := pgutils.NewClient()
	if err1 != nil {
		// print log
		log.Printf("[DATASTORE CONNECT] could not connect: %v", err1)
		// status = 1
		return utils.Result(1, "could not connect", "")
	}
	// status = 0
	return utils.Result(0, "", c)
}

func CreateSession(c pgutils.IDBClient, sessionId string) map[string]any {
	// try to create table if not exist
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})
	// insert data
	err1 := c.Insert(&pgutils.DeviceSession{
		SessionID:       sessionId,
		State:           "gwd:running|snmp:running",
		CreatedTime:     time.Now().String(),
		LastUpdatedTime: time.Now().String(),
	})
	if err1 != nil {
		// print log
		log.Printf("[DATASTORE CREATE] could not insert: %v", err1)
		// status = 1
		return utils.Result(1, "could not insert: "+err1.Error(), "")
	}
	// status = 0
	return utils.Result(0, "", sessionId)
}

func QuerySession(c pgutils.IDBClient, sessionId string) map[string]any {
	// try to create table if not exist
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})
	// define result
	var results = []pgutils.DeviceSession{}
	// query data
	err1 := c.Query(&results, pgutils.QueryExpr{
		Expr:  "session_id = ?",
		Value: sessionId,
	})
	if err1 != nil {
		// print log
		log.Printf("[DATASTORE QUERY STATUS] could not query: %v", err1)
		// status = 1
		return utils.Result(1, "could not query: "+err1.Error(), "")
	}
	if len(results) == 0 {
		// print log
		log.Printf("[DATASTORE QUERY STATUS] no data for this session ID ")
		// status = 2
		return utils.Result(2, "no data for this session ID", "")
	}
	if len(results) != 1 {
		// print log
		log.Printf("[DATASTORE QUERY STATUS] should only have one result")
		// status = 3
		return utils.Result(3, "should only have one result", "")
	}
	// status = 0
	return utils.Result(0, "", results[0])
}

func QueryResult(c pgutils.IDBClient, sessionId string) map[string]any {
	// try to create table if not exist
	c.CreateTable(&pgutils.DeviceResult{}, pgutils.CreateTableOpt{IfNotExists: true})
	// define result
	var results = []pgutils.DeviceResult{}
	// query data
	err1 := c.Query(&results, pgutils.QueryExpr{
		Expr:  "session_id = ?",
		Value: sessionId,
	})
	if err1 != nil {
		// print log
		log.Printf("[DATASTORE QUERY RESULT] could not query: %v", err1)
		// status = 1
		return utils.Result(1, "could not query: "+err1.Error(), "")
	}
	if len(results) == 0 {
		// print log
		log.Printf("[DATASTORE QUERY RESULT] no data for this session ID ")
		// status = 2
		return utils.Result(2, "no data for this session ID", "")
	}
	// status = 0
	return utils.Result(0, "", results)
}

func QueryLastSession(c pgutils.IDBClient) map[string]any {
	// try to create table if not exist
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})
	// define result1
	var results1 = pgutils.DeviceSession{}
	// define result2
	var results2 = []pgutils.DeviceSession{}
	// get database
	pgdb, err1 := c.GetDB()
	if err1 != nil {
		// print log
		log.Printf("[DATASTORE QUERY LAST SESSION] could not get database: %v", err1)
		// status = 1
		return utils.Result(1, "could not get database: "+err1.Error(), "")
	}
	// query max last_updated_time
	err2 := pgdb.Model((*pgutils.DeviceSession)(nil)).
		ColumnExpr("max(last_updated_time) as last_updated_time").
		Select(&results1)
	if err2 != nil {
		// print log
		log.Printf("[DATASTORE QUERY LAST SESSION] could not query max last_updated_time: %v", err2)
		// status = 2
		return utils.Result(2, "could not query max last_updated_time: "+err2.Error(), "")
	}
	if results1.LastUpdatedTime == "" {
		// print log
		log.Printf("[DATASTORE QUERY LAST SESSION] no data in device_sessions")
		// status = 3
		return utils.Result(3, "no data in device_sessions", "")
	}
	err4 := c.Query(&results2, pgutils.QueryExpr{
		Expr:  "last_updated_time = ?",
		Value: results1.LastUpdatedTime,
	})
	if err4 != nil {
		// print log
		log.Printf("[DATASTORE QUERY LAST SESSION] could not query: %v", err4)
		// status = 4
		return utils.Result(4, "could not query: "+err4.Error(), "")
	}
	if len(results2) != 1 {
		// print log
		log.Printf("[DATASTORE QUERY LAST SESSION] should only have one result")
		// status = 5
		return utils.Result(5, "should only have one result", "")
	}
	return utils.Result(0, "", results2[0])
}

func UpdateState(c pgutils.IDBClient, session pgutils.DeviceSession, state string) map[string]any {
	// try to create table if not exist
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})
	// get database
	pgdb, err1 := c.GetDB()
	if err1 != nil {
		// print log
		log.Printf("[DATASTORE UPDATE STATE] could not get database: %v", err1)
		// status = 1
		return utils.Result(1, "could not get database: "+err1.Error(), "")
	}
	// update data
	_, err2 := pgdb.Model(&pgutils.DeviceSession{
		ID:    session.ID,
		State: state,
	}).Column("state").WherePK().Update()
	if err2 != nil {
		// print log
		log.Printf("[DATASTORE UPDATE STATE] could not update state: %v", err2)
		// status = 2
		return utils.Result(2, "could not update state: "+err2.Error(), "")
	}
	return utils.Result(0, "", "")
}

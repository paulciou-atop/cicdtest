package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"nms/lib/pgutils"
	"nms/lib/pgutils/modules"
	"nms/lib/repo"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"

	lop "github.com/samber/lo/parallel"
)

var Topic = "scan.#"

type scanPayload struct {
	SessionID string `json:"sessionid,omitempty"`
}

// subscribe message
func ProcessScanResult(ctx context.Context, r repo.IRepo) {

	c := r.MQ()
	db := r.DB()
	rec := make(chan string)
	err := c.Subscribe(Topic, rec)
	if err != nil {
		logrus.Errorf("subscribe message queue fail: ", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			logrus.Info("Cancel subscribe")
			return
		case msg := <-rec:
			//receive scan result
			fmt.Println(msg)
			var receivedMsg scanPayload
			err := json.Unmarshal([]byte(msg), &receivedMsg)
			if err != nil {
				logrus.Errorf("Receive scan result %s decode fail: %v", msg, err)
				continue
			}
			err = updateInventories(db, receivedMsg.SessionID)
			if err != nil {
				logrus.Errorf("Process scan result %s  fail: %v", msg, err)
				continue
			}
		}
	}

}

// updateInventories update inventories from session result
func updateInventories(db pgutils.IClient, sessionId string) error {
	protocols, err := getScanProtocols(db, sessionId)
	if err != nil {
		return err
	}
	var devs []modules.DeviceResult
	db.Query(&devs, pgutils.QueryExpr{
		Expr:  "session_id = ?",
		Value: sessionId,
	})

	var devids []string
	for _, d := range devs {
		devids = append(devids, d.MacAddress)
	}

	missingInvs, err := findMissingInventories(db, devids)
	if err != nil {
		return err
	}
	// convert DeviceResult to Inventory
	marshalInvs := lop.Map(devs, func(d modules.DeviceResult, _ int) Inventory {
		return MarshalDeviceResultToInventory(d, protocols)
	})

	updateInvs := lop.Map(marshalInvs, func(inv Inventory, _ int) Inventory {
		return updateInventory(db, inv)
	})

	if len(missingInvs) > 0 {
		err = db.Update(&missingInvs)
		if err != nil {
			logrus.Errorf("Update missing inventories fail: ", err)
			return err
		}
	}

	if len(updateInvs) > 0 {
		pgClinet, err := db.GetDB()
		if err != nil {
			logrus.Errorf("get db client fail: ", err)
			return err
		}
		_, err = pgClinet.Model(&updateInvs).OnConflict("(id) DO UPDATE").Insert()
		if err != nil {
			logrus.Errorf("Update  inventories fail: ", err)
			return err
		}
	}

	return err
}

func getScanProtocols(db pgutils.IClient, sessionID string) ([]string, error) {
	// check session success
	var protocols []string
	var sessions []modules.DeviceSession
	db.Query(&sessions, pgutils.QueryExpr{
		Expr:  "session_id = ?",
		Value: sessionID,
	})
	if len(sessions) != 1 {
		return []string{}, ErrNotUnique(len(sessions))
	}
	sessionState := sessions[0].State
	exprs := strings.Split(sessionState, "|")
	for _, exp := range exprs {
		items := strings.Split(exp, ":")
		if len(items) != 2 {
			return []string{}, ErrSessionState
		}
		if strings.TrimSpace(items[1]) == "success" {
			protocols = append(protocols, items[0])
		}

	}
	return protocols, nil
}

func MarshalDeviceResultToInventory(dev modules.DeviceResult, protocols []string) Inventory {
	return Inventory{
		DeviceType: strings.Join(protocols, "|"),
		// Id will be assigned later, when check inv is a new or existing
		// Owner reserve for feature support
		Model: dev.Model,
		// Location reserve for feature support
		IpAddress:  dev.IpAddress,
		MacAddress: dev.MacAddress,
		HostName:   dev.Hostname,
		FirmwareInformation: FirmwareInfo{
			Ap:     dev.Ap,
			Kernel: dev.Kernel,
		},
		SupportProtocols: protocols,
	}
}

func findMissingInventories(db pgutils.IClient, devIDs []string) ([]Inventory, error) {
	var missingInvs []Inventory
	// Find missing divices
	pgdb, err := db.GetDB()
	if err != nil {
		return []Inventory{}, err
	}
	// Get missing inventory and update timestamp
	err = pgdb.Model(&missingInvs).
		Where("session_id not in (?)", pg.In(devIDs)).
		Select()
	missingInvs = lop.Map(missingInvs, func(inv Inventory, _ int) Inventory {
		return missingInv(inv)
	})
	return missingInvs, nil
}

// updateInventory update specific inventory, if inventory not exist add new one
func updateInventory(db pgutils.IClient, inv Inventory) Inventory {
	// TODO: we use mac addrss for now, change to real device id
	// when we upgrade device firmware and each device have unique ids
	currInv, err := findInventory(db, inv.MacAddress, "mac_address")
	if err != nil {
		//not found just add new one
		return newInv(inv)
	}
	updatedInv := updateInv(inv)
	err = mergo.MergeWithOverwrite(&currInv, updatedInv)
	if err != nil {
		logrus.Errorf("update inventory fail: %v", err)
	}
	return currInv
}

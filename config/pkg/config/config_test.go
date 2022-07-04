package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"nms/config/internal/services"
	"nms/config/pkg/session"
	"nms/lib/pgutils"
	"nms/lib/repo"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	metric1 := NewConfigMetric("test", "test", Metric{
		"ip":     "192.168.1.1",
		"port":   7766,
		"binary": []byte("abcde"),
	})

	metric2 := NewConfigMetric("test2", "test", Metric{
		"ip":     "192.168.1.1",
		"port":   7766,
		"binary": []byte("abcde"),
	})

	metric3 := NewConfigMetric("test2", "test", Metric{
		"ip":     "192.168.1.2",
		"port":   7766,
		"binary": []byte("abcde"),
	})

	metric4 := NewConfigMetric("test2", "test", Metric{
		"ip":     "192.168.1.2",
		"port":   7766,
		"binary": []byte("abde"),
	})

	if metric1.Hash() != metric2.Hash() {
		t.Errorf("Should be equal metric1.Hash() %s and metric2.Hash() %s", metric1.hash, metric2.hash)
	}

	if metric1.Hash() == metric3.Hash() {
		t.Errorf("Should not be equal metric1.Hash() %s and metric3.Hash() %s", metric1.hash, metric3.hash)
	}

	if metric1.Hash() == metric4.Hash() {
		t.Errorf("Should not be equal metric1.Hash() %s and metric4.Hash() %s", metric1.hash, metric4.hash)
	}

}

func TestStoreConfig(t *testing.T) {
	repo, err := repo.GetRepo(context.Background())
	if err != nil {
		t.Error(err)
	}

	db := repo.DB()
	pgdb, err := db.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS config_metrics`)
	if err != nil {
		t.Error(err)
	}

	db.CreateTable(&ConfigMetricModule{}, pgutils.CreateTableOpt{IfNotExists: true})
	// if err != nil {
	// 	t.Error("should not had error ", err)
	// 	return
	// }

	metric := NewConfigMetric("test", "test", Metric{
		"ip":     "192.168.1.1",
		"port":   7766,
		"binary": []byte("abcde"),
	})

	if err := StoreConfig(db, metric); err != nil {
		t.Error("should not have error ", err)
		return
	}

	var configs = []ConfigMetricModule{}
	if err := db.Query(&configs, pgutils.QueryExpr{
		Expr:  "hash = ?",
		Value: metric.hash,
	}); err != nil {
		t.Error("should not have error ", err)
		return
	}

	if len(configs) != 1 {
		t.Errorf("should only had 1 result but had %d", len(configs))
		return
	}

	payload := configs[0].Payload
	if payload["ip"].(string) != "192.168.1.1" {
		t.Errorf("config.payload[\"ip\"] should be 192.168.1.1 but %s", payload["ip"].(string))
	}
}

type Json = map[string]interface{}

func TestConfig(t *testing.T) {
	protocol := "dummyconfiger"
	kinds := []string{"general", "network", "snmp"}
	var devices []Device

	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	testDir := path.Join(dir, "dummydevs")
	os.Mkdir(testDir, os.ModePerm)

	// make sure db
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	r, err := repo.GetRepo(ctx)
	if err != nil {
		t.Error(err)
	}
	services.InitServices(r)
	InitTable(r.DB())
	session.InitDatabaseTables(r.DB())
	c := NewConfig(r)
	wg := new(sync.WaitGroup)
	wg.Add(5)
	var sessionIDs []string
	for i := 1; i < 6; i++ {
		filename := fmt.Sprintf("dev%d.json", i)
		filepath := path.Join(testDir, filename)
		ipaddress := fmt.Sprintf("10.0.12.1%d", i)
		macaddress := fmt.Sprintf("01:AB:CD:00:1E:%d0", i)
		emptydata := Json{}
		emptyJson, _ := json.MarshalIndent(emptydata, "", "  ")
		ioutil.WriteFile(filepath, emptyJson, 0644)

		// Prepare conifg metrics
		configTemp := map[string]Json{
			"general": {"name": filename},
			"network": {"ip": ipaddress, "mac": macaddress, "mask": "255.255.255.0"},
			"snmp":    {"port": 161, "private": "private", "ver": "V3"},
		}

		dev := Device{ID: filename, Path: filepath}
		devices = append(devices, dev)
		done := make(chan int)
		var metrics []*ConfigMetric
		for _, kind := range kinds {
			m := NewConfigMetric(protocol, kind, configTemp[kind])
			metrics = append(metrics, m)
		}

		go func() {
			ret, err := c.Config(ctx, dev, metrics, true, done)
			if err != nil {
				t.Error(err)
				goto Done
			}
			if ret.Session.State == "fail" {
				t.Errorf("session fail %s", ret.Device.DevicePath)
				goto Done
			}
			sessionIDs = append(sessionIDs, ret.Session.Id)
			<-done
		Done:
			wg.Done()

		}()

	}
	wg.Wait()
	//check result
	for _, sid := range sessionIDs {
		ret, err := session.GetConfigSession(r.DB(), sid)
		if err != nil {
			t.Error(err)
		}
		if ret.State != "success" {
			t.Errorf("session %s fail", ret.SessionID)
		}
	}
	// check config
	for _, dev := range devices {
		res, err := c.Upload(ctx, dev, protocol, []string{})
		if err != nil {
			t.Error(err)
		}
		if res.Success != true {
			t.Error(res.Message)
		}
		p := res.Payload.AsMap()
		gen := p["general"].(Json)
		name := gen["name"].(string)
		if name != dev.ID {
			t.Errorf("response name %s should be %s", name, dev.ID)
		}
	}

}

func TestCloseChan(t *testing.T) {
	c := make(chan int)
	rand.Seed(time.Now().UnixNano())
	// worker := func(out chan int) {
	// 	out <- rand.Intn(100)
	// }
	close(c)
	fmt.Println("out", <-c)
	fmt.Println("out", <-c)
	fmt.Println("out", <-c)
}

func TestTimeout(t *testing.T) {
	timeout := time.After(3 * time.Second)
	c := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 2)
		c <- struct{}{}
	}()
	select {
	case <-c:
		t.Log("receive c")
	case <-timeout:
		t.Error("time out")
	}
}

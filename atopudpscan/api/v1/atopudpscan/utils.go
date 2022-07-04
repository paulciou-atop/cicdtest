package atopudpscan

import (
	"bytes"
	context "context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nms/api/v1/serviceswatcher"
	"nms/atopudpscan/api/v1/dbstore"
	"os"
	"strconv"
	"strings"
	"time"

	"nms/api/v1/atopudpscan"
	mq "nms/messaging"

	"github.com/Atop-NMS-team/pgutils"
	"github.com/fatih/structs"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const datastore = "localhost:50051"

var postsql pgutils.IDBClient

func initPostgreSql() error {
	c, err := pgutils.NewClient()

	if err != nil {
		log.Print("can't connect pgutils datebase")
		return err
	}
	postsql = c
	return nil
}

func servicewatchercheck(servicename string) (string, error) {
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Can not connect to serviceswatcher")
		log.Printf("use default:%v\n", datastore)
		return datastore, err
	}
	defer conn.Close()

	client := serviceswatcher.NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// replace ServiceName to which you want to investigate
	res, err := client.Get(ctx, &serviceswatcher.GetRequest{
		ServiceName: servicename,
	})

	if err != nil {
		log.Println("Call servicewatcher.Watcher.Get err: ", err)
		return datastore, err
	}

	host := fmt.Sprintf("%s:%d", res.Info.Address, res.Info.Port)
	log.Printf("%v host=%s\n", servicename, host)

	return host, nil
}

func compare(a, b *atopudpscan.DeviceInfo) bool {
	a1, err := json.Marshal(a)
	if err != nil {

		return false
	}
	b1, err := json.Marshal(b)
	if err != nil {

		return false
	}
	r := bytes.Compare(a1, b1)
	if r == 0 {
		return true
	} else {
		return false
	}
}

//save data to DataStore
func saveToDataStore(id string, data []*atopudpscan.DeviceInfo) {
	l := len(data)
	if l <= 0 {
		log.Printf("saveToDatabase:no data")
		return
	}
	h, err := servicewatchercheck("datastore")
	if err != nil {
		return
	}
	log.Printf("datastore host:%v", h)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, h, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect to datebase: %v", err)
		return
	}

	dbInf_dbstore.TableName = name + "_" + id
	item := &dbstore.BlobData{Kind: "GWDScanResult", Size: int32(l), Element: make([]*dbstore.BlobData_BlobEle, 0)}

	for _, v := range data {
		m := structs.Map(v)
		i := 1
		e := &dbstore.BlobData_BlobEle{Kind: "row-" + strconv.Itoa(i), Size: int32(len(m)), KeyValueData: make([]*dbstore.BlobData_BlobEle_BlobKV, 0)}
		for k, v := range m {
			switch v := v.(type) {
			case string:
				e.KeyValueData = append(e.KeyValueData, &dbstore.BlobData_BlobEle_BlobKV{Key: k, Value: v})
			case bool:
				e.KeyValueData = append(e.KeyValueData, &dbstore.BlobData_BlobEle_BlobKV{Key: k, Value: strconv.FormatBool(v)})
			}
		}
		item.Element = append(item.Element, e)
	}

	c := dbstore.NewDbServiceClient(conn)
	_, err = c.DbWrite(context.Background(), &dbstore.DbWriteRequest{
		Bd:    item,
		DbInf: dbInf_dbstore,
	})
	if err != nil {
		log.Printf("error: %v", err)
	}

}

//save drvice to PostgreSql
func saveDeviceResult(s *Session, d *atopudpscan.DeviceInfo) {
	if postsql == nil {
		err := initPostgreSql()
		if err != nil {
			return
		}
	}

	postsql.CreateTable(&pgutils.DeviceResult{}, pgutils.CreateTableOpt{IfNotExists: true})

	data := &pgutils.DeviceResult{
		SessionID:  s.id,
		Model:      d.Model,
		MacAddress: d.MacAddress,
		IpAddress:  d.IPAddress,
		Netmask:    d.Netmask,
		Gateway:    d.Gateway,
		Hostname:   d.Hostname,
		Kernel:     d.Kernel,
		Ap:         d.Ap,
	}
	var result = []pgutils.DeviceResult{}
	err := postsql.Query(&result, pgutils.QueryExpr{Expr: "session_id = ? ", Value: s.id},
		pgutils.QueryExpr{Expr: "mac_address = ? ", Value: d.MacAddress})
	if err != nil {
		log.Print(err)
		return
	}

	if len(result) != 0 {
		for _, v := range result {
			data.ID = v.ID
			err := postsql.Update(data)
			if err != nil {
				log.Print(err)
			}
		}
	} else {
		err := postsql.Insert(data)
		if err != nil {
			log.Print(err)
		}
	}
}

//save session stats to PostgreSql
func saveDeviceSession(s *Session) {
	if postsql == nil {
		err := initPostgreSql()
		if err != nil {
			return
		}
	}

	postsql.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})

	var result = []pgutils.DeviceSession{}

	err := postsql.Query(&result, pgutils.QueryExpr{Expr: "session_id = ? ", Value: s.id})
	if err != nil {
		log.Print(err)
		return
	}
	if len(result) != 0 {
		for _, v := range result {
			data := &pgutils.DeviceSession{
				SessionID:       s.id,
				LastUpdatedTime: time.Now().String(),
				CreatedTime:     v.CreatedTime,
			}
			data.ID = v.ID
			data.State = updateSessionState(v.State, s.GetStatus().String())
			err := postsql.Update(data)
			if err != nil {
				log.Print(err)
			}
		}
	}

}
func updateSessionState(ori string, newstate string) string {
	data := map[string]string{}
	states := strings.Split(ori, "|")
	if len(states) != 2 {
		return newstate
	}
	for _, i := range states {
		kv := strings.Split(i, ":")
		if len(kv) != 2 {
			return newstate
		}
		data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	data["gwd"] = newstate
	return fmt.Sprintf("gwd:%s|snmp:%s", data["gwd"], data["snmp"])
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func publishDone(c mq.IMQClient, sessionid string) {
	go func() {
		mqc, err := mq.NewClient()
		if err != nil {
			log.Print(err)
		} else {
			msg := map[string]string{
				"sessionid": sessionid,
			}
			jsonret, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				return
			}
			mqc.Publish("scan.atopudpscan", string(jsonret))
		}
	}()
}

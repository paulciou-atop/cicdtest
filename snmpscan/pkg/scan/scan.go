/*
Package scan implements functions for scan devices which supports SNMP protocols.

*/
package scan

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"nms/snmpscan/pkg/snmp"
	"nms/snmpscan/pkg/store"
	"strconv"
	"sync"
	"time"

	mq "nms/messaging"

	"github.com/Atop-NMS-team/pgutils"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tatsushid/go-fastping"
)

// CIDR2IPs get ip list with given CIDR notation
func CIDR2IPs(cidr string) ([]net.IP, error) {
	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Errorf("")
		return []net.IP{}, err
	}

	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)
	var results []net.IP
	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		results = append(results, ip)

	}
	return results, nil
}

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

// GetReachableIPs
func GetReachableIPs(ips []net.IP) (results []net.IP) {

	p := fastping.NewPinger()

	pingResults := make(map[string]*response)
	for _, ip := range ips {

		target, err := net.ResolveIPAddr("ip4:icmp", ip.String())
		if err != nil {
			continue
		}

		pingResults[target.String()] = nil
		p.AddIPAddr(target)
	}

	onRecv, onIdle := make(chan *response), make(chan bool)
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	p.MaxRTT = time.Second
	p.RunLoop()

	for {
		select {
		case res := <-onRecv:
			if _, ok := pingResults[res.addr.String()]; ok {
				pingResults[res.addr.String()] = res
			}
		case <-onIdle:
			for _, r := range pingResults {
				if r != nil {
					results = append(results, r.addr.IP)
				}
			}
			return
		case <-p.Done():
			if err := p.Err(); err != nil {
				log.Error("ping fail : ", err)
			}
			return
		}
	}
}

type ScanResults = []snmp.ResultItem

func ScanAgent(cidr string, oids []string, opt *snmp.SnmpOption) (results []ScanResults, err error) {
	ips, err := CIDR2IPs(cidr)
	if err != nil {
		log.Errorf("CIDR %s format invalid : %+v\n", cidr, err)
		return
	}
	reachableDevs := GetReachableIPs(ips)
	fmt.Println("Reachable dev ", reachableDevs)
	var wg sync.WaitGroup
	for _, ip := range reachableDevs {
		wg.Add(1)
		go func(t net.IP) {
			ret, err := snmp.Get(t.String(), oids, opt)
			if err != nil {
				log.Infof("Snmp get %s err %v", t.String(), err)
				wg.Done()
				return
			}

			mac, err := snmp.FindFirstMac(t.String(), opt)
			if err == nil {
				ret = append(ret, snmp.ResultItem{
					Value: mac,
					Name:  "MAC",
					Kind:  "string",
				})
			}
			ret = append(ret, snmp.ResultItem{
				Value: t.String(),
				Name:  "IP",
				Oid:   "",
				Kind:  "string"})

			results = append(results, ret)

			wg.Done()
		}(ip)

	}
	wg.Wait()
	return
}

// getOid get object id value
func getOid(item ScanResults) string {

	for _, v := range item {

		if snmp.NormalizedOid(v.Name) == ".1.3.6.1.2.1.1.2.0" {
			s, ok := v.Value.(string)
			if ok {
				return snmp.NormalizedOid(s)
			}
		}
	}
	return ""
}

//ScanAtopDevices scan CIDR subnet and list all snmp agents which were producted by Atop.
func ScanAtopDevices(cidr string, opts *snmp.SnmpOption) (results []ScanResults, err error) {
	ips, err := CIDR2IPs(cidr)
	if err != nil {
		log.Errorf("CIDR %s format invalid : %+v\n", cidr, err)
		return
	}

	reachableDevs := GetReachableIPs(ips)
	fmt.Println("Reachable dev ", reachableDevs)
	var wg sync.WaitGroup
	for _, ip := range reachableDevs {
		wg.Add(1)
		go func(t net.IP) {
			ret, oid, err := snmp.GetAtopSystemInfo(t.String(), opts)
			if err != nil {
				//log.Infof("Snmp get %s atop system infor  err %v", t.String(), err)
				wg.Done()
				return
			}

			// Add IP
			ret = append(ret, snmp.ResultItem{
				Name:  "IP",
				Kind:  "string",
				Value: t.String(),
				Oid:   "",
			})

			atopSupportList := viper.GetStringSlice("atopDevices")
			for _, id := range atopSupportList {
				// filter devices which is Atop's device
				if snmp.NormalizedOid(oid) == snmp.NormalizedOid(id) {
					results = append(results, ret)
				}
			}

			wg.Done()
		}(ip)

	}
	wg.Wait()

	return
}

///////////////////////////////////////////////////
// Asynchronous functions

const category = "ss"
const path = "/snmpscan/asyncscan"

/////////////////////////////////////
// state handy tools

type ScanStoreStruct struct {
	State     string        `bson:"state" json:"state" structs:"state" mapstructure:"state"`
	TimeStamp string        `bson:"timestamp" json:"timestamp" structs:"timestamp" mapstructure:"timestamp"`
	Error     string        `bson:"error" json:"error" structs:"error" mapstructure:"error"`
	Data      []ScanResults `bson:"data" json:"data" structs:"data" mapstructure:"data"`
}

const RFC3339 = "2006-01-02T15:04:05Z07:00"

// VToString inferface{} to string
func VToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(RFC3339)
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%T:%v", v, v)
	}
}

func ScanStoreStructToTable(sessionid string, in ScanStoreStruct) (ScanStoreTable, []ScanResultRow) {
	scanResultFKey := sessionid
	table := ScanStoreTable{
		Key:            sessionid,
		SessionID:      sessionid,
		State:          in.State,
		TimeStamp:      in.TimeStamp,
		Error:          in.Error,
		ScanResultFKey: scanResultFKey,
	}
	// ScanResult Table
	scanResultTable := []ScanResultRow{}

	for _, v := range in.Data {
		for _, row := range v {
			scanResultTable = append(scanResultTable, ScanResultRow{
				Key:   scanResultFKey,
				Value: VToString(row.Value),
				Kind:  row.Kind,
				Name:  row.Name,
				Oid:   row.Oid,
			})
		}
	}
	return table, scanResultTable
}

type ScanStoreTable struct {
	Key            string `bson:"key" json:"key" structs:"key" mapstructure:"key"`
	SessionID      string `bson:"sessionID" json:"sessionID" structs:"sessionID" mapstructure:"sessionID"`
	State          string `bson:"state" json:"state" structs:"state" mapstructure:"state"`
	TimeStamp      string `bson:"timestamp" json:"timestamp" structs:"timestamp" mapstructure:"timestamp"`
	Error          string `bson:"error" json:"error" structs:"error" mapstructure:"error"`
	ScanResultFKey string `bson:"scanResultKey" json:"scanResultKey" structs:"scanResultKey" mapstructure:"scanResultKey"`
}

type ScanResultRow struct {
	Key   string `bson:"key" json:"key" structs:"key" mapstructure:"key"`
	Value string `bson:"value" json:"value" structs:"value" mapstructure:"value"`
	Name  string `bson:"name" json:"name" structs:"name" mapstructure:"name"`
	Kind  string `bson:"kind" json:"kind" structs:"kind" mapstructure:"kind"`
	Oid   string `bson:"oid" json:"oid" structs:"oid" mapstructure:"oid"`
}

var initState = ScanStoreStruct{
	State:     "running",
	TimeStamp: time.Now().String(),
	Data:      []ScanResults{},
}

func errState(err error) ScanStoreStruct {
	return ScanStoreStruct{Error: err.Error(), TimeStamp: time.Now().String(), State: "fail"}
}

func finalState(data []ScanResults) ScanStoreStruct {
	return ScanStoreStruct{
		State:     "success",
		TimeStamp: time.Now().String(),
		Data:      data,
	}
}

//CreateSession greate session ID
func CreateSession() string {
	id := store.CreateID(category)
	store.Store(path, id, structs.Map(initState))
	return id
}

func SessionFail(id string, err error) {
	d := errState(err)
	_, retErr := store.Store(path, id, structs.Map(&d))
	if retErr != nil {
		log.Error("store session fail state fail: ", retErr)
	}

	if retErr != nil {
		log.Error("postgres store session fail state fail: ", retErr)
	}
}

type JsonObj = map[string]interface{}

var ctx context.Context

// AsyncScan asynchronous scan
func AsyncScan(id string, cidr string, oids []string, opts *snmp.SnmpOption, justAtop bool) {

	var results []ScanResults
	ss, err := store.NewSession(id)
	ss.StoreSession((&pgutils.DeviceSession{
		State: "running",
	}))

	ss.StoreSession((&pgutils.DeviceSession{
		State: "running",
	}))

	if justAtop {
		results, err = ScanAtopDevices(cidr, opts)
	} else {
		results, err = ScanAgent(cidr, oids, opts)
	}

	if err != nil {
		ss.StoreSession(&pgutils.DeviceSession{
			State: "fail",
		})
		return
	}
	data := finalState(results)
	m := structs.Map(&data)
	devResults := []pgutils.DeviceResult{}
	for _, dev := range results {
		ret := pgutils.DeviceResult{}
		for _, kv := range dev {
			value := fmt.Sprintf("%v", kv.Value)
			switch kv.Name {
			case "SystemDescr":
				ret.Description = value
			case "SystemFwVer":
				ret.FirmwareVer = value
			case "SystemMacAddress":
				ret.MacAddress = value
			case "SystemKernelVer":
				ret.Kernel = value
			case "SystemModel":
				ret.Model = value
			case "IP":
				ret.IpAddress = value
			}
		}
		ret.SessionID = id
		devResults = append(devResults, ret)
	}
	if len(devResults) <= 0 {
		ss.StoreSession(&pgutils.DeviceSession{State: "notfound"})
	} else {
		ss.StoreSession(&pgutils.DeviceSession{State: "success"})
		ss.StoreDeviceResults(devResults)
	}

	store.Store(path, id, m)

	go func() {
		mqc, err := mq.NewClient()
		if err != nil {
			log.Error(err)
		} else {
			msg := map[string]string{
				"sessionid": id,
			}
			jsonret, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				return
			}
			mqc.Publish("scan.snmpscan", string(jsonret))
		}
		defer mqc.Close()
	}()

	ss.Close()
}

type AsyncScanResult struct {
	SessionID string        `bson:"session-id" json:"session-id" structs:"session-id"`
	State     string        `bson:"state" json:"state" structs:"state"`
	TimeStamp string        `bson:"timestamp" json:"timestamp" structs:"timestamp" mapstructure:"timestamp"`
	Payload   []ScanResults `bson:"payload" json:"payload" structs:"payload"`
}

func GetAsyncResult(seesionID string) (AsyncScanResult, error) {
	r, err := store.Read(path, seesionID)
	if err != nil {
		log.Errorf("store read id %s has error %+v", seesionID, err)
		return AsyncScanResult{}, err
	}

	var structData ScanStoreStruct
	err = mapstructure.Decode(r.Payload, &structData)
	if err != nil {
		log.Errorf("convert result of read session %s has error %+v", seesionID, err)
		return AsyncScanResult{}, err
	}
	state, ok := r.Payload["state"].(string)
	if !ok {
		return AsyncScanResult{}, fmt.Errorf("data convert fail in session %s's result ", seesionID)
	}
	return AsyncScanResult{
		SessionID: seesionID,
		TimeStamp: structData.TimeStamp,
		State:     state,
		Payload:   structData.Data,
	}, nil
}

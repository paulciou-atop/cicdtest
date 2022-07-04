/*
Package snmp implements bunch of functions for snmp get/set information from targets

Type convert
	typeconvert.go have bunch of snmp response covert functions
*/

package snmp

import (
	"fmt"
	"time"

	g "github.com/gosnmp/gosnmp"
	glog "github.com/sirupsen/logrus"
)

// Target "192.168.1.10"

func NormalizedOid(oid string) string {
	if oid[0] != '.' {
		return "." + oid
	}
	return oid
}

// findFirstMac find the first mac address in the target
func FindFirstMac(target string, opt *SnmpOption) (mac string, err error) {

	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err = client.Connect()
	if err != nil {
		glog.Errorf("Connect() err: %v", err)
		return
	}
	defer client.Conn.Close()

	pdus, err := client.WalkAll(ifPhysAddress)

	if err != nil {
		glog.Errorf("WalkAll %s err: %v", target, err)
		return
	}
	for _, v := range pdus {
		if v.Type == g.OctetString {
			if len(v.Value.([]byte)) > 0 {

				return convertoctec(v), nil
			}
		}
	}
	return "", fmt.Errorf("can't found MAC address on  %s", target)
}

// getSubTree walk all children below rootOid
func getSubTree(target, rootOid string, opt *SnmpOption) (results []g.SnmpPDU, err error) {
	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err = client.Connect()
	if err != nil {
		glog.Errorf("Connect() err: %v", err)
		return
	}
	defer client.Conn.Close()

	results, err = client.WalkAll(rootOid)

	if err != nil {
		glog.Errorf("WalkAll() with %v err: %v ", opt, err)
		return
	}
	return
}

//Walk
func Walk(target, rootOid string, opt *SnmpOption) (results []ResultItem, err error) {
	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err = client.Connect()
	if err != nil {
		glog.Errorf("Connect() err: %v", err)
		return
	}
	defer client.Conn.Close()

	pdus, err := client.WalkAll(rootOid)
	if err != nil {
		glog.Error("WalkAll() err: %v", err)
		return
	}
	results = convertSnmpPDUsToResultItems(pdus, map[string]string{})
	return
}

type ResultItem struct {
	Value interface{} `json:"value" structs:"value" mapstructure:"value"`
	Name  string      `json:"name" structs:"name" mapstructure:"name"`
	Kind  string      `json:"kind" structs:"kind" mapstructure:"kind"`
	Oid   string      `bson:"oid" json:"oid" structs:"oid" mapstructure:"oid"`
}

type SnmpOption struct {
	// Port is a port.
	Port uint16
	// Community is an SNMP Community string.
	Community string
	// Version is an SNMP Version.
	Version g.SnmpVersion
	// Timeout is the timeout for one SNMP request/response.
	Timeout time.Duration
	// Set the number of retries to attempt.
	Retries int
}

var DefaultSnmpOption = SnmpOption{
	Port:      161,
	Community: "public",
	Version:   g.Version2c,
	Timeout:   time.Second,
	Retries:   1,
}

func getOids(target string, oids []string, opt *SnmpOption) (results []g.SnmpPDU, err error) {
	if opt == nil {
		opt = &DefaultSnmpOption
	}

	if len(oids) == 0 {
		oids = BasicOIDs
	}

	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err = client.Connect()
	if err != nil {
		glog.Errorf("Connect() err: %v", err)
		return
	}
	defer client.Conn.Close()
	pkt, err := client.Get(oids)
	if err != nil {
		//glog.Errorf("snmp get %s err: %v", target, err)
		return
	}
	return pkt.Variables, err

}

//Get target is a IP address or host url which we want to get oids' value, opt = nil means use default settings
func Get(target string, oids []string, opt *SnmpOption) (results []ResultItem, err error) {

	pdus, err := getOids(target, oids, opt)

	if err != nil {
		return
	}
	results = convertSnmpPDUsToResultItems(pdus, map[string]string{})

	return results, nil
}

//Bulk get same type below oid
func Bulk(target string, oids []string, opt *SnmpOption) (results []ResultItem, err error) {
	if len(oids) == 0 {
		oids = BasicOIDs
	}

	if opt == nil {
		opt = &DefaultSnmpOption
	}
	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err = client.Connect()
	if err != nil {
		glog.Errorf("Connect() err: %v", err)
		return
	}
	defer client.Conn.Close()
	pdus, err := client.GetBulk(oids, 0, 50)
	if err != nil {
		return
	}
	for _, v := range pdus.Variables {
		fmt.Printf("pdus %v\n", v)
	}

	results = convertSnmpPDUsToResultItems(pdus.Variables, map[string]string{})

	return results, nil
}

// Describe
func Describe(target string, opt *SnmpOption) (PublicMatrix, error) {
	if opt == nil {
		opt = &DefaultSnmpOption
	}
	return getPublicMatrix(target, opt)
}

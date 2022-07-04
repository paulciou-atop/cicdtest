// This file consist of some common utility functions
package cmd

import (
	"bufio"
	"nms/snmpscan/pkg/snmp"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
	glog "github.com/sirupsen/logrus"
)

var versionMap = map[string]uint8{
	"v1": 0x0,
	"v2": 0x1,
	"v3": 0x3,
}

// converSNMPVer convert version stirng from CLI to gosnmp verson
func converSNMPVer(in string) uint8 {
	v, ok := versionMap[in]
	if !ok {
		return 0x01
	}
	return v
}

//loadOIDsFile loading OIDs file and convert to string slice (OIDs)
func loadOIDsFile(path string) ([]string, error) {
	var oids []string
	reader, err := os.Open(path)
	if err != nil {
		glog.Error("OIDs file not found, use default oids instead of ")
		return oids, err
	}

	bufReader := bufio.NewReader(reader)
	lines, _, err := bufReader.ReadLine()
	if err == nil {
		for _, line := range lines {
			oids = append(oids, string(line))
		}
	}

	return oids, nil
}

type SnmpOptions struct {
	Port          int32
	ReadCommunity string
	Version       string
	Timeout       int32
	Retries       int32
}

func SnmpOptions2SnmpOption(in *SnmpOptions) *snmp.SnmpOption {
	opt := &snmp.SnmpOption{
		Port:      uint16(in.Port),
		Community: in.ReadCommunity,
		Version:   gosnmp.SnmpVersion(converSNMPVer(in.Version)),
		Timeout:   time.Duration(in.Timeout) * time.Second,
		Retries:   int(in.Retries),
	}
	return opt
}

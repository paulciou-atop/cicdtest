package snmp

import (
	"fmt"
	"reflect"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

/*
This file consist of atop private oids
*/

var systemInfoOids = []SnmpUnitItem{
	{
		Name: "SystemDescr",
		Type: String,
		Oid:  ".1.4.0",
	},
	{
		Name: "SystemFwVer",
		Type: String,
		Oid:  ".1.5.0",
	},
	{
		Name: "SystemMacAddress",
		Type: Mac,
		Oid:  ".1.6.0",
	},
	{
		Name: "SystemKernelVer",
		Type: String,
		Oid:  ".1.7.0",
	},

	{
		Name: "SystemModel",
		Type: String,
		Oid:  ".1.10.0",
	},
}

// AtopSystemInfoMatrix
type AtopSystemInfoMatrix struct {
	SystemDescr      string `bson:"systemDescr" json:"systemDescr" structs:"systemDescr" mapstructure:"systemDescr"`
	SystemFwVer      string `bson:"systemFwVer" json:"systemFwVer" structs:"systemFwVer" mapstructure:"systemFwVer"`
	SystemMacAddress string `bson:"systemMacAddress" json:"systemMacAddress" structs:"systemMacAddress" mapstructure:"systemMacAddress"`
	SystemKernelVer  string `bson:"systemKernelVer" json:"systemKernelVer" structs:"systemKernelVer" mapstructure:"systemKernelVer"`
	SystemModel      string `bson:"systemModel" json:"systemModel" structs:"systemModel" mapstructure:"systemModel"`
}

func setStringToAtopSystemInfoMatrix(items []g.SnmpPDU, matrix *AtopSystemInfoMatrix, field string, convertfunc func(interface{}) string) {
	for _, v := range items {
		strValue := convertfunc(v.Value)
		if strValue != "" {
			reflect.ValueOf(matrix).Elem().FieldByName(field).SetString(strValue)
		}
	}
}

// GetAtopSystemInfo get atop device system informaiton
func GetAtopSystemInfo(target string, opt *SnmpOption) (results []ResultItem, objectID string, err error) {
	snmpRes, err := getOids(target, []string{SystemObjectID}, opt)
	if err != nil {
		// log.Errorf("snmp get value from %s err: %+v", target, err)
		return
	}

	if len(snmpRes) != 1 {
		err = fmt.Errorf("did not get object id from %s", target)
		log.Error(err)
		return
	}
	objectID = toString(snmpRes[0].Value)
	fmt.Printf("object %v(%T) id is %s\n", snmpRes[0].Value, snmpRes[0].Value, objectID)
	oids := []string{}
	nameMapping := map[string]string{}
	for _, v := range systemInfoOids {
		oids = append(oids, objectID+v.Oid)
		nameMapping[objectID+v.Oid] = v.Name
	}

	pdus, err := getOids(target, oids, opt)
	if err != nil {
		log.Info("snmp get value from %s err: %+v", target, err)
		return
	}

	results = convertSnmpPDUsToResultItems(pdus, nameMapping)
	return
}

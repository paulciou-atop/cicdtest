package snmp

import (
	"math/big"
	"reflect"
	"strings"

	g "github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
)

//////////////////////////////////////////////////
// Interface section

type ItemType int

const (
	String ItemType = iota
	Mac
	Number
)

// IntefaceItem describe item data module in the interface
type SnmpUnitItem struct {
	Name string
	Type ItemType
	Oid  string
}

const ifIndex = ".1.3.6.1.2.1.2.2.1.1"
const ifPhysAddress = ".1.3.6.1.2.1.2.2.1.6"

var interfacesOids = []SnmpUnitItem{

	{
		Name: "Descr",
		Type: String,
		Oid:  ".1.3.6.1.2.1.2.2.1.2",
	},
	{
		Name: "PhysAddress",
		Type: Mac,
		Oid:  ".1.3.6.1.2.1.2.2.1.6",
	},
	{
		Name: "AdminStatus", // set/get port enable/disable
		Type: Number,
		Oid:  ".1.3.6.1.2.1.2.2.1.7",
	},
	{
		Name: "OperStatus", // get port status
		Type: Number,
		Oid:  ".1.3.6.1.2.1.2.2.1.8",
	},
	{
		Name: "InNUcastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.2.2.1.12",
	},
	{
		Name: "InErrors",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.2.2.1.14",
	}, {
		Name: "OutNucastPkts", // ifOutMulticastPkts + ifOutBroadcastPkt
		Type: Number,
		Oid:  ".1.3.6.1.2.1.2.2.1.18",
	},
	{
		Name: "OutErrors",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.2.2.1.20",
	},
	{
		Name: "HCInOctets",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.6",
	}, {
		Name: "HCInUcastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.7",
	}, {
		Name: "HCInMulticastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.8",
	}, {
		Name: "HCInBroadcastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.9",
	}, {
		Name: "HCOutOctets",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.10",
	}, {
		Name: "HCOutUcastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.11",
	},
	{
		Name: "HCOutMulticastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.12",
	}, {
		Name: "HCOutBroadcastPkts",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.13",
	}, {
		Name: "HighSpeed", // get port spee
		Type: Number,
		Oid:  ".1.3.6.1.2.1.31.1.1.1.15",
	},
}

type InterfaceMatrix struct {
	Index              string `bson:"index" json:"index" structs:"index"`
	Descr              string `bson:"description" json:"description" structs:"description"`
	PhysAddress        string `bson:"MAC" json:"MAC" structs:"MAC"`
	AdminStatus        int32  `bson:"portEnable" json:"portEnable" structs:"portEnable"`
	OperStatus         int32  `bson:"portStatus" json:"portStatus" structs:"portStatus"`
	InNUcastPkts       int64  `bson:"inNUcastPkts" json:"inNUcastPkts" structs:"inNUcastPkts"`
	InErrors           int64  `bson:"inErrors" json:"inErrors" structs:"inErrors"`
	OutNucastPkts      int64  `bson:"outNUcastPkts" json:"outNUcastPkts" structs:"outNUcastPkts"`
	OutErrors          int64  `bson:"outErrors" json:"outErrors" structs:"outErrors"`
	HCInOctets         int64  `bson:"inOctets" json:"inOctets" structs:"inOctets"`
	HCInUcastPkts      int64  `bson:"InUcastPkts" json:"InUcastPkts" structs:"InUcastPkts"`
	HCInMulticastPkts  int64  `bson:"inMulticastPkts" json:"inMulticastPkts" structs:"inMulticastPkts"`
	HCInBroadcastPkts  int64  `bson:"inBroadcastPkts" json:"inBroadcastPkts" structs:"inBroadcastPkts"`
	HCOutOctets        int64  `bson:"outOctets" json:"outOctets" structs:"outOctets"`
	HCOutUcastPkts     int64  `bson:"outUcastPkts" json:"outUcastPkts" structs:"outUcastPkts"`
	HCOutMulticastPkts int64  `bson:"outMulticastPkts" json:"outMulticastPkts" structs:"outMulticastPkts"`
	HCOutBroadcastPkts int64  `bson:"outBroadcastPkts" json:"outBroadcastPkts" structs:"outBroadcastPkts"`
	HighSpeed          int32  `bson:"highSpeed" json:"highSpeed" structs:"highSpeed"`
}

func setStringToInterfaceMatrix(items []g.SnmpPDU, matrix *[]InterfaceMatrix, field string, convertfunc func(interface{}) string) {
	for _, v := range items {
		strValue := convertfunc(v.Value)
		if strValue != "" {
			for i, idx := range *matrix {
				if strings.HasSuffix(v.Name, idx.Index) {
					reflect.ValueOf(&(*matrix)[i]).Elem().FieldByName(field).SetString(strValue)
				}
			}
		}
	}
}

func setNumberToInterfaceMatrix(items []g.SnmpPDU, matrix *[]InterfaceMatrix, field string, convertfunc func(interface{}) *big.Int) {
	for _, v := range items {
		intValue := convertfunc(v.Value)

		for i, idx := range *matrix {
			if strings.HasSuffix(v.Name, idx.Index) {
				structfield := reflect.ValueOf(&(*matrix)[i]).Elem().FieldByName(field)
				if structfield != (reflect.Value{}) {
					structfield.SetInt(intValue.Int64())
				} else {
					logrus.Errorf("interface struct field not exist: %s", field)
				}

			}
		}

	}
}

// GetInterface get interface matrix
func GetInterface(target string, opt *SnmpOption) (results []InterfaceMatrix, err error) {

	ifIndexs, err := getSubTree(target, ifIndex, opt)
	if err != nil {
		logrus.Error("snmp get sub-tree err: %+v", err)
		return
	}

	//create index
	for _, v := range ifIndexs {
		idx := toIndexString(v.Value)
		if idx != "" {
			results = append(results, InterfaceMatrix{Index: idx})
		}
	}

	for _, v := range interfacesOids {
		pdus, err := getSubTree(target, v.Oid, opt)
		if err != nil {
			logrus.Error("snmp get sub-tree err: %+v", err)
			continue
		}
		switch v.Type {
		case Number:
			setNumberToInterfaceMatrix(pdus, &results, v.Name, toBigInt)
		case Mac:
			setStringToInterfaceMatrix(pdus, &results, v.Name, toMacString)
		case String:
			setStringToInterfaceMatrix(pdus, &results, v.Name, toString)
		}
	}

	return
}

//////////////////////////////////////////////////
// Public Matrix

var publicOids = []SnmpUnitItem{
	{
		Name: "SysDescr",
		Type: String,
		Oid:  ".1.3.6.1.2.1.1.1.0",
	},
	{
		Name: "SysObjectId",
		Type: String,
		Oid:  ".1.3.6.1.2.1.1.2.0",
	}, {
		Name: "SysUpTime",
		Type: Number,
		Oid:  ".1.3.6.1.2.1.1.3.0",
	}, {
		Name: "SysContact",
		Type: String,
		Oid:  ".1.3.6.1.2.1.1.4.0",
	}, {
		Name: "SysName",
		Type: String,
		Oid:  ".1.3.6.1.2.1.1.5.0",
	}, {
		Name: "SysLocation",
		Type: String,
		Oid:  ".1.3.6.1.2.1.1.6.0",
	}, {
		Name: "SysServicesDescr",
		Type: String,
		Oid:  ".1.3.6.1.2.1.1.7.0",
	},
}

type PublicMatrix struct {
	SysDescr    string            `protobuf:"bytes,1,opt,name=sysDescr,proto3" json:"sysDescr,omitempty"`
	SysObjectId string            `protobuf:"bytes,2,opt,name=sysObjectId,proto3" json:"sysObjectId,omitempty"`
	SysUpTime   int64             `protobuf:"varint,3,opt,name=sysUpTime,proto3" json:"sysUpTime,omitempty"`
	SysContact  string            `protobuf:"bytes,4,opt,name=sysContact,proto3" json:"sysContact,omitempty"`
	SysName     string            `protobuf:"bytes,5,opt,name=sysName,proto3" json:"sysName,omitempty"`
	SysLocation string            `protobuf:"bytes,6,opt,name=sysLocation,proto3" json:"sysLocation,omitempty"`
	SysServices int32             `protobuf:"varint,7,opt,name=sysServices,proto3" json:"sysServices,omitempty"`
	Interfaces  []InterfaceMatrix `protobuf:"bytes,8,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
}

func setStringToPublicMatrix(items []g.SnmpPDU, matrix *PublicMatrix, item SnmpUnitItem, convertfunc func(interface{}) string) {
	for _, v := range items {

		if NormalizedOid(v.Name) == NormalizedOid(item.Oid) {
			strValue := convertfunc(v.Value)
			if strValue != "" {
				structfield := reflect.ValueOf(matrix).Elem().FieldByName(item.Name)
				if structfield != (reflect.Value{}) {
					structfield.SetString(strValue)
				} else {
					logrus.Errorf("interface struct field not exist: %s", item.Name)
				}
			}
		}

	}
}

func setNumberToPublicMatrix(items []g.SnmpPDU, matrix *PublicMatrix, item SnmpUnitItem, convertfunc func(interface{}) *big.Int) {
	for _, v := range items {

		if NormalizedOid(v.Name) == NormalizedOid(item.Oid) {
			intValue := convertfunc(v.Value)

			structfield := reflect.ValueOf(matrix).Elem().FieldByName(item.Name)
			if structfield != (reflect.Value{}) {
				structfield.SetInt(intValue.Int64())
			} else {
				logrus.Errorf("interface struct field not exist: %s", item.Name)
			}
		}

	}
}

// getPublicMatrix get target's public matrix
func getPublicMatrix(target string, opt *SnmpOption) (publicMatrix PublicMatrix, err error) {
	oids := []string{}
	for _, v := range publicOids {
		oids = append(oids, v.Oid)
	}
	pdus, err := getOids(target, oids, opt)
	for _, v := range publicOids {
		switch v.Type {
		case Number:
			setNumberToPublicMatrix(pdus, &publicMatrix, v, toBigInt)
		case Mac:
			setStringToPublicMatrix(pdus, &publicMatrix, v, toMacString)
		case String:
			setStringToPublicMatrix(pdus, &publicMatrix, v, toString)
		}
	}

	interfaces, err := GetInterface(target, opt)
	if err != nil {
		logrus.Errorf("snmpscan get interface err : %+v", err)
		return publicMatrix, err
	}
	publicMatrix.Interfaces = interfaces
	return
}

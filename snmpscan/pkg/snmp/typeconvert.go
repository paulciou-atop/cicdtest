package snmp

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

func convertSnmpPDUsToResultItems(in []g.SnmpPDU, nameMapping map[string]string) []ResultItem {
	var results []ResultItem
	for _, v := range in {
		var t interface{}
		var kind string
		switch v.Type {
		case g.OctetString:
			t = convertoctec(v)
			kind = "string"
		case g.ObjectIdentifier:
			t = v.Value.(string)
			kind = "string"
		default:
			t = g.ToBigInt(v.Value)
			kind = "int"
		}
		name, ok := nameMapping[v.Name]
		if !ok {
			name = v.Name
		}
		results = append(results, ResultItem{
			Value: t,
			Oid:   v.Name,
			Name:  name,
			Kind:  kind,
		})
	}
	return results
}

func insertdash(s string) string {
	upper := strings.ToUpper(s)
	var buffer bytes.Buffer
	var n_1 = 2 - 1
	var l_1 = len(upper) - 1
	for i, rune := range upper {
		buffer.WriteRune(rune)
		if i%2 == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}

func convertoctec(in g.SnmpPDU) string {

	if in.Type != g.OctetString {
		return ""
	}
	value := string(in.Value.([]byte))

	//Mac
	if strings.HasPrefix(NormalizedOid(in.Name), NormalizedOid(ifPhysAddress)) {
		return toMacString(in.Value)
	}
	// Atop Mac
	data := in.Value.([]byte)
	if len(data) == 6 && data[0] == 0x00 && data[1] == 0x60 {
		return toMacString(in.Value)
	}
	return value
}

// toString to normal string fail return empty string ""
func toString(value interface{}) string {
	bytes, ok := value.([]byte)
	if ok {
		return string(bytes)
	}
	s, ok := value.(string)
	if ok {
		return s
	}
	return ""
}

// toMacString to mac address string fail return empty string ""
func toMacString(value interface{}) string {
	bytes, ok := value.([]byte)
	if ok {
		temp := hex.EncodeToString(bytes)
		return insertdash(temp)
	}
	return ""
}

// toIndexString number to 10-base string, if fail return empty string ""
func toIndexString(value interface{}) string {
	v, ok := value.(int)
	if ok {
		str := strconv.Itoa(v)
		return str
	}
	return ""
}

func toBigInt(value interface{}) *big.Int {
	return g.ToBigInt(value)
}

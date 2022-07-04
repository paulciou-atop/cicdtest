package store_test

import (
	"fmt"
	"nms/snmpscan/pkg/store"
	"strings"
	"testing"
)

func TestGetSchema(t *testing.T) {
	m := map[string]interface{}{
		"word":   "string",
		"number": 1234,
		"float":  12.334,
		"array":  []string{"1", "2"},
		"map": map[string]interface{}{
			"map.string": "abc",
		},
	}
	ret := store.GetShema(m)
	t.Log(ret)
	if ret["array"] != "[]string" {
		t.Error("array should be []string")
	}
	if ret["float"] != "float64" {
		t.Error("float should be float64")
	}
	if ret["word"] != "string" {
		t.Error("word should be string")
	}
	t.Log(ret)

	//t.Log(store.GetShema(m["map"].(map[string]interface{})))
}

func updateState(ori string, newstate string) string {
	data := map[string]string{}
	states := strings.Split(ori, "|")

	for _, i := range states {
		kv := strings.Split(i, ":")
		if len(kv) < 2 {
			data["gwd"] = ""
		} else {
			data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	data["snmp"] = newstate
	return fmt.Sprintf("gwd:%s|snmp:%s", data["gwd"], data["snmp"])
}

func TestUpdatState(t *testing.T) {
	case1 := "snmp:123"
	newcase1 := updateState(case1, "222")
	if newcase1 != "gwd:|snmp:222" {
		t.Error("should be gwd:|snmp:222 but ", newcase1)
	}

	case2 := "snmp:123|gwd:4444"
	newcase2 := updateState(case2, "222")
	if newcase2 != "gwd:4444|snmp:222" {
		t.Error("should be gwd:4444|snmp:222 but ", newcase2)
	}

	case3 := "gwd : 4444 | snmp : 123"
	newcase3 := updateState(case3, "222")
	if newcase3 != "gwd:4444|snmp:222" {
		t.Error("should be gwd:4444|snmp:222 but ", newcase3)
	}

	case4 := "gwd:123"
	newcase4 := updateState(case4, "222")
	if newcase4 != "gwd:123|snmp:222" {
		t.Error("should be gwd:123|snmp:222 but ", newcase4)
	}

	case5 := "running"
	newcase5 := updateState(case5, "222")
	if newcase5 != "gwd:|snmp:222" {
		t.Error("should be gwd:|snmp:222 but ", newcase5)
	}

}

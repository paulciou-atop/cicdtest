package device

import "testing"

var map1 = map[string]interface{}{
	"a": "a",
	"b": 2,
}

var map2 = map[string]interface{}{
	"b": "4",
	"d": "5",
}

func TestMapMerge(t *testing.T) {
	m := mergeMap(map1, map2)
	t.Log(m)
}

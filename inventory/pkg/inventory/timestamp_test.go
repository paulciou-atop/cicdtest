package inventory

import "testing"

func appendA(inv Inventory) Inventory {
	r := inv
	r.Name += "A"
	return r
}

func appendB(inv Inventory) Inventory {
	r := inv
	r.Name += "B"
	return r
}

func TestCompose(t *testing.T) {
	inv := Inventory{
		Name: "inv",
	}
	ret := pipe(appendA, appendB, appendB, appendA)(inv)
	if ret.Name != "invABBA" {
		t.Errorf("name should be invABBA but %s", ret.Name)
	}
}

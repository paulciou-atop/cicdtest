package inventory

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"nms/lib/pgutils"
)

var ErrNotUnique = func(l int) error {
	return fmt.Errorf("should be unique but get %d results", l)
}
var ErrSessionState = errors.New("scan session state error")

func sha1CheckSum(data any) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	h.Write(jsonBytes)

	return hex.EncodeToString(h.Sum(nil)), nil
}

func findInventory(db pgutils.IClient, deviceID, idColumn string) (Inventory, error) {
	var invs []Inventory
	expr := fmt.Sprintf("%s = ?", idColumn)
	err := db.Query(&invs, pgutils.QueryExpr{
		Expr:  expr,
		Value: deviceID,
	})
	if err != nil {
		return Inventory{}, err
	}
	if len(invs) != 1 {
		return Inventory{}, ErrNotUnique(len(invs))
	}
	return invs[0], nil
}

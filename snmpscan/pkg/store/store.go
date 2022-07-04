/*
Package stroe implements bunch of functions for data store, even storeservice or any other internal cache
*/
package store

import (
	"fmt"
	"reflect"
	// "snmpscan/api/v1/dataStore"

	"github.com/google/uuid"
)

type JsonObj = map[string]interface{}

// Result result of functions
type Result struct {
	Path    string  `bson:"path" json:"path" structs:"path" mapstructure:"path"`
	Id      string  `bson:"id" json:"id" structs:"id" mapstructure:"id"`
	Payload JsonObj `bson:"payload" json:"payload" structs:"payload" mapstructure:"payload"`
}

// CreateID create unique id for record. for example snmpscan's result
// creates ss:77152046-a853-4b6b-82f3-40920427ca12
func CreateID(category string) string {
	return category + ":" + uuid.New().String()
}

// GetColumns get record's columns
func GetColumns(record JsonObj) []string {
	result := []string{}
	for key, _ := range record {
		result = append(result, key)
	}
	return result
}

// GetShema
// result's type actually is a map[string]string, key is column name, value is type
func GetShema(record JsonObj) map[string]string {
	if record == nil {
		return map[string]string{}
	}
	result := map[string]string{}
	for i, v := range record {
		t := reflect.TypeOf(v)

		result[i] = t.String()
	}
	return result
}

// Storer interface, function to store data
type Storer interface {
	// Store store data into path with specific id
	// id: unique id, could be a session id, identify
	// path: where store data, example: /db/scanresult /cache/scanresult /db/timeseries/scan/matrix
	Store(path string, id string, data JsonObj) (Result, error)
}

// Reader interface, function to read data from abstract store
type Reader interface {
	Read(path, id string) (Result, error)
}

// Client interface, store client
type Client interface {
	Storer
	Reader
}

// func StoreToService(db, dbtype, tablename string, data []JsonObj) error {
// 	return dataStore.Write(db, dbtype, tablename, data)
// }

func Store(path string, id string, data JsonObj) (Result, error) {

	//TODO: Add storeservice here to store result
	rc := RedisClient()
	if rc == nil {
		return Result{}, fmt.Errorf("got a nil redis client")
	}
	return rc.Store(path, id, data)
}

func Read(path, id string) (Result, error) {
	//TODO: Add storeservice here to store result
	rc := RedisClient()
	if rc == nil {
		return Result{}, fmt.Errorf("got a nil redis client")
	}
	return rc.Read(path, id)
}

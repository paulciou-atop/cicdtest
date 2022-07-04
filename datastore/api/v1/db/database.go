package db

import (
	"context"
	"datastore/api/v1/proto"
	db_store "datastore/api/v1/proto/dataStore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	iniDbMapping()
}

//var collection *mongo.Collection

type Server struct {
	statusStore.SessionStatusServiceServer
}
type Server2 struct {
	db_store.DbServiceServer
}
type BlobData struct {
	Kind    string //db name
	Size    int32  //?
	Element []*BlobEle
}

type WriteToDb struct {
	Kind string
	Gzz  string
}

type BlobEle struct {
	Kind         string //table name
	Size         int32  //how many columns
	KeyValueData []*BlobKV
}

type BlobKV struct {
	Key string
	//Value *anypb.Any
	Value string
}

type IDataBase interface {
	CreateOne(dbInf DbInfo, ctx context.Context, document interface{}) (oidHex string, err error)
	DbIni() *mongo.Client //temp ver , try to change to how to return no specific dbclient
	//====below has not implement
	ListDb() error
	ReadDb() error
	CreatMany() error
	UpdateDb() error
	DeleteRecord() error
}
type MongoDB struct {
}
type MySql struct {
}

var dbMapping = make(map[string]interface{})

type ssItem struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Ip     string             `bson:"ip"`
	Status string             `bson:"status"`
	Type   string             `bson:"type"`
	Time   string             `bson:"time"`
	Msg    string             `bson:"msg"`
}
type DbInfo struct {
	DbType string
	DbName string
	TbName string
}

func iniDbMapping() {
	dbMapping["mongo"] = func() struct{} {
		return MongoDB{}
	}
	dbMapping["mysql"] = func() struct{} {
		return MySql{}
	}
}

func saveOneRecord(I IDataBase, dbInf DbInfo, ctx context.Context, document interface{}) (oidHex string, err error) {
	return I.CreateOne(dbInf, ctx, document)
}

func initialDatabase(I IDataBase) *mongo.Client {

	return I.DbIni()
}

func listDBStream(I IDataBase, ) error {
	return I.ListDb()
}
func readDbByCondition(I IDataBase, ) error {
	return I.ReadDb()
}
func createMany(I IDataBase, ) error {
	return I.CreatMany()
}
func updateDb(I IDataBase, ) error {
	return I.UpdateDb()
}
func deleteDb(I IDataBase, ) error {
	return I.DeleteRecord()
}



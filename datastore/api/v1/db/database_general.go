package db

// This file  stores data in db by BlobData,and will modify to general purpose version.
import (
	"context"
	db_store "datastore/api/v1/proto/dataStore"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (*Server2) DbWrite(ctx context.Context, req *db_store.DbWriteRequest) (*db_store.DbWriteResponse, error) {
	// save data in db
	mDB := MongoDB{}

	log.Println("call Dbwrite()333")
	bd := req.GetBd()
	//dsItem := db_store.BlobData{
	//	Kind: "bbb",
	//	Size: 5,
	//	Element: []*db_store.BlobData_BlobEle{
	//		{
	//			Kind: "a-2",
	//			Size: 2,
	//			KeyValueData: []*db_store.BlobData_BlobEle_BlobKV{
	//				{Key: "c_3",
	//					Value: "317",
	//				},
	//				{Key: "c_4",
	//					Value: "449"},
	//			},
	//		},
	//	},
	//}

	//can work
	data := BlobData{
		Kind:    bd.GetKind(),
		Size:    bd.GetSize(),
		Element: dataStoreBlobeEleToBlobEle(bd.Element),
	}
	log.Println(data)

	dbInfo := db_storeDbInfTrans(req.GetDbInf())
	oidHex, _ := saveOneRecord(mDB, dbInfo, ctx, data)
	log.Println(oidHex)

	//fmt.Println(data)
	return &db_store.DbWriteResponse{
		Bd: &db_store.BlobData{
			Kind:    bd.GetKind(),
			Size:    bd.GetSize(),
			Element: bd.Element,
		},
	}, nil

}

func (*Server2) ListDb(_ *db_store.ListDbRequest, stream db_store.DbService_ListDbServer) error {
	//list all device status from db
	mDB := MongoDB{}
	client := initialDatabase(mDB)
	collection := client.Database("nms").Collection("device_data_store")
	defer DbClose(client)
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background()) // Should handle err
	for cur.Next(context.Background()) {
		data := &BlobData{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)

		}
		stream.Send(&db_store.ListDbResponse{Bd: dataTransformBlobData(data)}) // Should handle err
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

//---private funcion
func dataStoreBlobKVToBlobKV(blobkv []*db_store.BlobData_BlobEle_BlobKV) []*BlobKV {
	var returnArray []*BlobKV
	for _, item := range blobkv {

		temp := BlobKV{}
		temp.Key = item.Key
		temp.Value = item.Value

		returnArray = append(returnArray, &temp)
	}
	return returnArray

}

func blobKVTodataStoreBlobKV(blobkv []*BlobKV) []*db_store.BlobData_BlobEle_BlobKV {
	var returnArray []*db_store.BlobData_BlobEle_BlobKV
	for _, item := range blobkv {

		temp := db_store.BlobData_BlobEle_BlobKV{}
		temp.Key = item.Key
		temp.Value = item.Value

		returnArray = append(returnArray, &temp)
	}
	return returnArray

}

func db_storeDbInfTrans(req *db_store.DatabaseInfo) DbInfo {
	return DbInfo{

		DbType: req.GetDatabaseType(),
		DbName: req.GetDatabaseName(),
		TbName: req.GetTableName(),
	}
}

func dataTransformBlobEle(blobEle []*BlobEle) []*db_store.BlobData_BlobEle {
	var returnArray []*db_store.BlobData_BlobEle
	for _, item := range blobEle {

		temp := db_store.BlobData_BlobEle{}
		temp.Kind = item.Kind
		temp.Size = item.Size
		temp.KeyValueData = blobKVTodataStoreBlobKV(item.KeyValueData)
		returnArray = append(returnArray, &temp)
	}
	return returnArray

}

func dataStoreBlobEle(blobEle []*db_store.BlobData_BlobEle) []*BlobEle {
	var returnArray []*BlobEle
	for _, item := range blobEle {

		temp := BlobEle{}
		temp.Kind = item.Kind
		temp.Size = item.Size
		temp.KeyValueData = dataStoreBlobKVToBlobKV(item.KeyValueData)
		returnArray = append(returnArray, &temp)
	}
	return returnArray

}

func dataStoreBlobeEleToBlobEle(blobEle []*db_store.BlobData_BlobEle) []*BlobEle {
	var returnArray []*BlobEle
	for _, item := range blobEle {

		temp := BlobEle{}
		temp.Kind = item.Kind
		temp.Size = item.Size
		temp.KeyValueData = dataStoreBlobKVToBlobKV(item.KeyValueData)
		returnArray = append(returnArray, &temp)
	}
	return returnArray

}

func dataTransformBlobData(data *BlobData) *db_store.BlobData {

	ele := dataTransformBlobEle(data.Element)
	return &db_store.BlobData{
		Kind:    data.Kind,
		Size:    data.Size,
		Element: ele,
	}

}

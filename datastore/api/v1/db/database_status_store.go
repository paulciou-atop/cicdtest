package db

//old version and will maintain  until the general version has done
import (
	"context"
	statusStore "datastore/api/v1/proto"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func reqDbInfTrans(req *statusStore.DatabaseInfo) DbInfo {
	return DbInfo{

		DbType: req.GetDatabaseType(),
		DbName: req.GetDatabaseName(),
		TbName: req.GetTableName(),
	}
}

func (*Server) CreateSessionStatus(ctx context.Context, req *statusStore.CreateSessionStatusRequest) (*statusStore.CreateSessionStatusResponse, error) {
	// save data in db
	log.Println("create on record in createSessionStatus func")

	ss := req.GetSs()
	//reqDbType := req.GetDbInfo().GetDatabaseType()
	dbInfo := reqDbInfTrans(req.GetDbInfo())

	data := ssItem{
		Ip:     ss.GetIp(),
		Status: ss.GetStatus(),
		Type:   ss.GetType(),
		Time:   ss.GetTime(),
		Msg:    ss.GetMsg(),
	}
	mDB := MongoDB{}
	//dbInterface := dbMapping[reqDbType]
	oidHex, _ := saveOneRecord(mDB, dbInfo, ctx, data)

	return &statusStore.CreateSessionStatusResponse{
		Ss: &statusStore.SessionStatus{
			Id:     oidHex,
			Ip:     ss.GetIp(),
			Status: ss.GetStatus(),
			Type:   ss.GetType(),
			Msg:    ss.GetMsg(),
		},
	}, nil

}

func (*Server) GetSessionStatus(ctx context.Context, req *statusStore.GetSessionStatusRequest) (*statusStore.GetSessionStatusResponse, error) {
	//query data by id
	mDB := MongoDB{}
	client := initialDatabase(mDB)
	dbInfo := reqDbInfTrans(req.GetDbInfo())

	collection := client.Database(dbInfo.DbName).Collection(dbInfo.TbName)
	ssId := req.GetId()

	oid, err := primitive.ObjectIDFromHex(ssId)
	// create an empty struct
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
	data := &ssItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find record with specified IP: %v", err),
		)
	}

	return &statusStore.GetSessionStatusResponse{
		Ss: dataTostatusStore(data),
	}, nil
}

func dataTostatusStore(data *ssItem) *statusStore.SessionStatus {
	return &statusStore.SessionStatus{
		Id:     data.ID.Hex(),
		Ip:     data.Ip,
		Status: data.Status,
		Type:   data.Type,
		Time:   data.Time,
		Msg:    data.Msg,
	}
}

//func dataToBlobStore(data *BlobData) *db_store.BlobData {
//	return &db_store.BlobData{
//		Kind:    data.Kind,
//		Size:    data.Size,
//		Element: *data.Element,
//	}
//
//}

//update by IP
func (*Server) UpdateSessionStatus(ctx context.Context, req *statusStore.UpdateSessionStatusRequest) (*statusStore.UpdateSessionStatusResponse, error) {
	fmt.Println("Update start")
	ss := req.GetSs()
	oid, err := primitive.ObjectIDFromHex(ss.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("--------Cannot parse ID-------"),
		)
	}

	mDB := MongoDB{}

	client := initialDatabase(mDB)

	collection := client.Database("nms").Collection("device_data_store")
	defer DbClose(client)
	// create an empty struct
	data := &ssItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find record with specified id: %v", err),
		)
	}

	// we update our internal struct

	data.Ip = ss.GetIp()
	data.Time = ss.GetTime()
	data.Status = ss.GetStatus()
	data.Type = ss.GetType()
	data.Msg = ss.GetMsg()
	log.Println("the data we want to update", data)
	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	return &statusStore.UpdateSessionStatusResponse{
		Ss: dataTostatusStore(data),
	}, nil

}

func (*Server) ListSessionStatus(_ *statusStore.ListSessionStatusRequest, stream statusStore.SessionStatusService_ListSessionStatusServer) error {
	//list all device status from db
	mDB := MongoDB{}

	client := initialDatabase(mDB)
	collection := client.Database("nms").Collection("device_data_store")
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background()) // Should handle err
	for cur.Next(context.Background()) {
		data := &ssItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)

		}
		stream.Send(&statusStore.ListSessionStatusResponse{Ss: dataTostatusStore(data)}) // Should handle err
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

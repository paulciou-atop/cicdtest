package main

import (
	"context"
	db_store "datastore/api/v1/proto/dataStore"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"

	"datastore/api/v1/proto"
	"google.golang.org/grpc"
)

func testCreateRecord(cc grpc.ClientConnInterface, dbInf *statusStore.DatabaseInfo, req *statusStore.SessionStatus) (dsRes *statusStore.CreateSessionStatusResponse, err error) {

	//dsItem := &statusStore.SessionStatus{
	//	Ip:     "192.168.0.122",
	//	Status: "status1",
	//	Type:   "snmp",
	//	Msg:    "newMsg",
	//	Time:   "2022/04/18 20:22:33",
	//}

	c := statusStore.NewSessionStatusServiceClient(cc)
	createSsRes, err := c.CreateSessionStatus(context.Background(), &statusStore.CreateSessionStatusRequest{
		Ss:     req,
		DbInfo: dbInf,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	return createSsRes, err
}
func testQueryById(cc grpc.ClientConnInterface, queryId string) (dsRes *statusStore.GetSessionStatusResponse, err error) {
	fmt.Println("get  ds data by Ip")

	c := statusStore.NewSessionStatusServiceClient(cc)
	_, err2 := c.GetSessionStatus(context.Background(), &statusStore.GetSessionStatusRequest{Id: queryId})
	if err2 != nil {
		fmt.Printf("Error happened while getting data: %v \n", err2)
	}

	getSsReq := &statusStore.GetSessionStatusRequest{Id: queryId}
	getSsRes, getSsErr := c.GetSessionStatus(context.Background(), getSsReq)
	if getSsErr != nil {
		fmt.Printf("Error happened while reading: %v \n", getSsErr)
	}

	fmt.Printf("device status: %v \n", getSsRes)
	return getSsRes, getSsErr
}
func testUpdateService(cc grpc.ClientConnInterface, req *statusStore.SessionStatus) (dsRes *statusStore.UpdateSessionStatusResponse, err error) {
	c := statusStore.NewSessionStatusServiceClient(cc)
	//fmt.Println("get  ds data by Ip")

	//res, err := testQueryByIp(cc, Ip)
	//if err != nil {
	//	fmt.Println("Error happened while querying ip: %v \n ", err)
	//}
	//recordId := res.GetSs().GetId()

	dsItem := &statusStore.SessionStatus{
		Id:     req.GetId(),
		Ip:     req.GetIp(),
		Time:   req.GetTime(),
		Status: req.GetStatus(),
		Type:   req.GetType(),
		Msg:    req.GetMsg(),
	}
	updateRes, updateErr := c.UpdateSessionStatus(context.Background(), &statusStore.UpdateSessionStatusRequest{Ss: dsItem})

	return updateRes, updateErr
}
func testDBWrite(cc grpc.ClientConnInterface, dbInf *db_store.DatabaseInfo) {

	dsItem := db_store.BlobData{
		Kind: "data_store_table",
		Size: 5,
		Element: []*db_store.BlobData_BlobEle{
			{
				Kind: "row-1",
				Size: 30,
				KeyValueData: []*db_store.BlobData_BlobEle_BlobKV{
					{Key: "c_3",
						Value: "317",
					},
					{Key: "c_4",
						Value: "449"},
				},
			},
			//{
			//	Kind: "row-2",
			//	Size: 44,
			//	KeyValueData: []*db_store.BlobData_BlobEle_BlobKV{
			//		{Key: "c_3",
			//			Value: "317",
			//		},
			//		{Key: "c_4",
			//			Value: "449"},
			//	},
			//},
		},
	}
	c := db_store.NewDbServiceClient(cc)
	createSsRes, err := c.DbWrite(context.Background(), &db_store.DbWriteRequest{
		Bd:    &dsItem,
		DbInf: dbInf,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Println(createSsRes)
}

func testListDB(cc grpc.ClientConnInterface) {
	c := db_store.NewDbServiceClient(cc)
	stream, err := c.ListDb(context.Background(), &db_store.ListDbRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBd())
	}
}

func runCreateSessionStore(dbInf *statusStore.DatabaseInfo, cc grpc.ClientConnInterface) {
	ssItem := &statusStore.SessionStatus{
		Ip:     "192.168.0.144",
		Status: "status13333",
		Type:   "snmp",
		Msg:    "newMsg",
		Time:   "2022/04/18 20:22:33",
	}
	createSsRes, err := testCreateRecord(cc, dbInf, ssItem)
	//create Ss

	//createSsRes, err := c.CreateSessionStatus(context.Background(), &statusStore.CreateSessionStatusRequest{Ss: dsItem})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("device status has been created: %v", createSsRes)
}
func runQueryById(dbInf *statusStore.DatabaseInfo, cc grpc.ClientConnInterface) {
	fmt.Println("get  ds data by Id")
	c := statusStore.NewSessionStatusServiceClient(cc)
	queryId := "625c8425efea55b10bd5e567"
	getSsReq := &statusStore.GetSessionStatusRequest{Id: queryId,
		DbInfo: dbInf,
	}
	getSsRes, getSsErr := c.GetSessionStatus(context.Background(), getSsReq)
	if getSsErr != nil {
		fmt.Printf("Error happened while reading: %v \n", getSsErr)
	}

	fmt.Printf("device status: %v \n", getSsRes)
}
func runUpdateById(dbInf *statusStore.DatabaseInfo, cc grpc.ClientConnInterface) {
	fmt.Println("update data by Id")
	//c := statusStore.NewSessionStatusServiceClient(cc)
	queryId := "625c8425efea55b10bd5e567"
	ssItem := &statusStore.SessionStatus{
		Id:     queryId,
		Ip:     "192.168.0.126",
		Status: "status1",
		Type:   "snmp",
		Msg:    "newMsg",
		Time:   "2022/04/18 20:22:33",
	}

	updateRes, updateErr := testUpdateService(cc, ssItem)
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf(" was updated: %v\n", updateRes)
}
func runTestListdb(dbInf *statusStore.DatabaseInfo, cc grpc.ClientConnInterface) {
	//list all records
	c := statusStore.NewSessionStatusServiceClient(cc)
	stream, err := c.ListSessionStatus(context.Background(), &statusStore.ListSessionStatusRequest{
		DbInfo: dbInf,
	})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetSs())
	}
}
func main() {

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	dbInf := &statusStore.DatabaseInfo{
		DatabaseType: "mongo",
		DatabaseName: "nms",
		TableName:    "device_data_store",
	}
	//dbInf_dbstore := &db_store.DatabaseInfo{
	//	DatabaseType: "mongo",
	//	DatabaseName: "nms",
	//	TableName:    "device_data_store",
	//}
	//
	runCreateSessionStore(dbInf, cc)
	runQueryById(dbInf, cc)
	runUpdateById(dbInf, cc)
	runTestListdb(dbInf, cc)
	
	println("===========DB Write=========")
	//testDBWrite(cc, dbInf_dbstore)
	//println("==========LIST DB=============")
	//testListDB(cc)
}

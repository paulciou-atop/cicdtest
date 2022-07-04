package scan

import (
	"context"
	"fmt"
	"io"
	"testing/api/v1/serviceswatcher"
	"time"

	"github.com/bobbae/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Scan(cidr string) (err error) {
	conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return
	}
	defer conn.Close()

	client := serviceswatcher.NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// replace ServiceName to which you want to investigate
	res, err := client.Get(ctx, &serviceswatcher.GetRequest{
		ServiceName: "snmpscan",
	})

	if err != nil {
		glog.Error("Call servicewatcher.Watcher.Get err: ", err)
		return
	}

	host := fmt.Sprintf("%s:%d", res.Info.Address, res.Info.Port)

	connScan, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return
	}
	defer conn.Close()
	scanClient := NewSnmpScanClient(connScan)
	// Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (SnmpScan_ScanClient, error)
	stream, err := scanClient.StartAsyncScan(context.Background())
	if err != nil {
		fmt.Printf("Scan err:%+v\n", err)
		return
	}
	fmt.Println("Send start async request")
	stream.Send(&ScanRequest{
		Range: cidr,
		SnmpSettings: &SnmpSettings{
			Port:           161,
			ReadCommunity:  "public",
			WriteCommunity: "private",
			Version:        SnmpSettings_ver2c},
		Oids: []string{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&ScanRequest{
		Range: cidr,
		SnmpSettings: &SnmpSettings{
			Port:           161,
			ReadCommunity:  "public",
			WriteCommunity: "private",
			Version:        SnmpSettings_ver2c},
		Oids: []string{},
	})

	for {

		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error while reading stream: %v\n", err)
			continue
		}
		fmt.Printf("Receive sessionid from StartAsyncScan: %v\n", msg.SessionId)
		fmt.Println("Waiting 10 sec for scan")
		time.Sleep(time.Second * 30)
		asyncResult, err := scanClient.GetAsyncScanResult(context.Background(), &AsyncRequest{
			SessionId: msg.SessionId,
		})
		if err != nil {
			fmt.Printf("error while reading stream: %v\n", err)
			continue
		}
		fmt.Printf("session %s state:%s\n", msg.SessionId, asyncResult.Status)
		for i, singleResult := range asyncResult.Result {
			fmt.Printf("%dth device\n", i)
			for _, unit := range singleResult.Pdus {
				fmt.Printf(">>>>>>>>> %s:%s\n", unit.Name, unit.Value)
			}

		}
	}
	return
}

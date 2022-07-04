package service

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	up "nms/api/v1/atopudpscan"
	ww "nms/api/v1/serviceswatcher"

	"nms/scanservice/api/utils"
)

func GetGwdHost() map[string]any {
	// open client
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err1 := grpc.Dial("serviceswatcher:8081", opts)
	if err1 != nil {
		// print log
		log.Printf("[SERVICESWATCHER] could not connect: %v", err1)
		// status = 1
		return utils.Result(1, "could not connect", "")
	}
	defer conn.Close()
	// connect to serviceswatcher
	watcher := ww.NewWatcherClient(conn)
	// get snmpscan host
	host, err2 := watcher.Get(context.Background(), &ww.GetRequest{ServiceName: "atopudpscan"})
	if err2 != nil {
		// print log
		log.Printf("[SERVICESWATCHER] call servicewatcher.Watcher.Get error: %v", err2)
		// status = 2
		return utils.Result(2, "call servicewatcher.Watcher.Get error", "")
	}
	// print log
	log.Printf("address: %s; port: %d", host.Info.Address, host.Info.Port)
	// status = 0
	return utils.Result(0, "", fmt.Sprintf("%s:%d", host.Info.Address, host.Info.Port))
}

func ConnectGWD(host string) map[string]any {
	// connect to gwdscan
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err1 := grpc.Dial(host, opts)
	if err1 != nil {
		// print log
		log.Printf("[GWDSCAN START] could not connect: %v", err1)
		// status = 1
		return utils.Result(1, "could not connect", "")
	}
	c := up.NewGwdClient(conn)
	// status = 0
	return utils.Result(0, "", c)
}

func GwdStartScan(c up.GwdClient, ServerIp string, sessionId string) map[string]any {
	// start scanning
	res, err1 := c.SessionScan(context.Background(), &up.ScanConfig{
		ServerIp: ServerIp,
		Id:       sessionId,
	})
	if err1 != nil {
		log.Printf("[GWDSCAN START] can not scan %v", err1)
		// status = 1
		return utils.Result(1, "can not scan: "+err1.Error(), "")
	}
	// print response
	log.Printf("[START GWD]: %v \n", res)
	// has error message or not
	if res.GetMessage() != "" {
		// status = 2
		return utils.Result(2, res.GetMessage(), "")
	}
	// status = 0
	return utils.Result(0, "", res.GetId())
}

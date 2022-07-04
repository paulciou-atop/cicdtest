package service

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ss "nms/api/v1/snmpscan"
	pb "nms/api/v1/scanservice"
	ww "nms/api/v1/serviceswatcher"

	"nms/scanservice/api/utils"
)

func GetSnmpHost() map[string]any {
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
	host, err2 := watcher.Get(context.Background(), &ww.GetRequest{ServiceName: "snmpscan"})
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

func ConnectSNMP(host string) map[string]any {
	// connect to snmpscan
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err1 := grpc.Dial(host, opts)
	if err1 != nil {
		// print log
		log.Printf("[SNMPSCAN START] could not connect: %v", err1)
		// status = 1
		return utils.Result(1, "could not connect", "")
	}
	c := ss.NewSnmpScanClient(conn)
	// status = 0
	return utils.Result(0, "", c)
}

func SnmpStartScan(c ss.SnmpScanClient, Range string, SnmpSettings *pb.SnmpSettings, Oids []string, sessionId string) map[string]any {
	// convert &pb.SnmpSettings to &ss.SnmpSettings
	scanRequest := &ss.ScanRequest{
		Range:       Range,
		AtopDevices: true,
		Oids:        Oids,
	}
	if SnmpSettings != nil {
		scanRequest.SnmpSettings = &ss.SnmpSettings{
			Port:           SnmpSettings.Port,
			ReadCommunity:  SnmpSettings.ReadCommunity,
			WriteCommunity: SnmpSettings.WriteCommunity,
			Version:        ss.SnmpSettings_SNMPVer(SnmpSettings.Version),
		}
	}
	// start scanning
	res, err1 := c.SessionScan(context.Background(), &ss.SessionScanRequest{
		SessionId:   sessionId,
		ScanRequest: scanRequest,
	})
	if err1 != nil {
		// print log
		log.Printf("[SNMPSCAN START] can not scan %v", err1)
		// status = 1
		return utils.Result(1, "can not scan: "+err1.Error(), "")
	}
	// print response
	log.Printf("[START SNMP]: %v \n", res)
	// status = 0
	return utils.Result(0, "", "")
}

package scan

import (
	"context"
	"exporter/api/v1/serviceswatcher"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DeviceInfo struct {
	IP     string
	UpTime int64
}

func Scan(cidr string) (devinfos []DeviceInfo, err error) {
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
		log.Error("Call servicewatcher.Watcher.Get err: ", err)
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
	stream, err := scanClient.Scan(context.Background(), &ScanRequest{
		Range:       cidr,
		AtopDevices: true,
	})
	if err != nil {
		log.Error("snmp.Scan err ", err)
		return
	}

	for {

		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error while reading stream: %v\n", err)
			continue
		}

		for _, v := range msg.Pdus {
			if v.Name == "IP" {
				// ips = append(ips, v.Value.GetStringValue())
				ip := v.Value.GetStringValue()
				ret, err := scanClient.Describe(context.Background(), &DescribeRequest{
					Target: ip,
				})
				if err != nil {
					log.Infof("Describe %s error err %v", ip, err)
					continue
				}
				devinfos = append(devinfos, DeviceInfo{
					IP:     ip,
					UpTime: ret.Information.SysUpTime,
				})
			}
		}

	}
	return
}

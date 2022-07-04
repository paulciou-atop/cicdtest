package snmp_test

import (
	"context"
	"fmt"
	swAPI "nms/api/v1/serviceswatcher"
	"nms/snmpscan/pkg/scan"
	"nms/snmpscan/pkg/store"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestScan(t *testing.T) {
	ips, err := scan.CIDR2IPs("192.168.13.55/24")
	if err != nil {
		t.Error("GetAllIps err:", err)
	}
	t.Log("len(ips) ", len(ips))
	assert.Equal(t, 256, len(ips))
}

func getServiceURL(name string) (string, error) {
	conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return "", err
	}
	defer conn.Close()
	client := swAPI.NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	res, err := client.Get(ctx, &swAPI.GetRequest{
		ServiceName: "snmpscan",
	})

	if err != nil {
		log.Error("Call servicewatcher.Watcher.Get err: ", err)
		return "", err
	}
	return fmt.Sprintf("%s:%d", res.Info.Address, res.Info.Port), nil
}

func TestRedis(t *testing.T) {
	err := store.ExampleClient()
	if err != nil {
		t.Error("Redis example error ", err)
	}
}

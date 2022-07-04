package serviceswatcher

import (
	"fmt"
	"time"

	api "nms/api/v1/serviceswatcher"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ServiceName = "snmpscan"

func Register(addr string, port int32, kind []string) error {
	conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Infoln("Can not connect to serviceswatcher:8081 try to connect localhost.")
		conn, err = grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error("Can not connect to servicewatcher niether conatainer nor localhost")
			return err
		}
	}
	defer conn.Close()
	client := api.NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err = client.Register(ctx, &api.ServiceInfo{
		Name:    ServiceName,
		Address: addr,
		Port:    port,
		Kind:    kind,
	})

	if err != nil {
		log.Error("Call servicewatcher.Register err: ", err)
		return err
	}

	return nil

}

// GetServiceHost
func GetServiceHost(name string) (string, error) {
	conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return "", err
	}
	defer conn.Close()
	client := api.NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	res, err := client.Get(ctx, &api.GetRequest{
		ServiceName: name,
	})

	if err != nil {
		log.Error("Call servicewatcher.Watcher.Get err: ", err)
		return "", err
	}
	return fmt.Sprintf("%s:%d", res.Info.Address, res.Info.Port), nil
}

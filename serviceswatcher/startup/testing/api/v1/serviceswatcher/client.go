package serviceswatcher

import (
	"context"
	"fmt"
	"time"

	glog "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Testing() {
	conn, err := grpc.Dial("serviceswatcher:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Error("Can not connect to serviceswatcher")
		return
	}
	defer conn.Close()

	client := NewWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// replace ServiceName to which you want to investigate
	res, err := client.Get(ctx, &GetRequest{
		ServiceName: "mongo",
	})

	if err != nil {
		glog.Error("Call servicewatcher.Watcher.Get err: ", err)
		return
	}

	host := fmt.Sprintf("%s:%d", res.Info.Address, res.Info.Port)
	fmt.Printf("mongo host = %s\n", host)

}

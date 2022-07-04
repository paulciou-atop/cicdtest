package main

import (
	"datastore/api/v1/db"
	"datastore/api/v1/proto"
	db_store "datastore/api/v1/proto/dataStore"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//client := db.DbIni("mongoDB", "nms", "device_data_store")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	db_store.RegisterDbServiceServer(s, &db.Server2{})
	statusStore.RegisterSessionStatusServiceServer(s, &db.Server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	//db.DbClose(client)
	s.Stop()
	fmt.Println("End of Program")
}

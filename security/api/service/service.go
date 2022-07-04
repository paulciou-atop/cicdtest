package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"security/api/v1/security"
	"security/config"

	"github.com/bobbae/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var grpcport = config.GetgrpcPort()
var httpport = config.GethttpPort()

func main() {
	lis, err := net.Listen("tcp", ":"+grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	s := security.NewPkiService()
	security.RegisterPkiServer(grpcServer, s)
	reflection.Register(grpcServer)
	log.Printf("Start service.....%v\n", grpcport)
	glog.V(2).Infoln("security Serving gRPC on 0.0.0.0", grpcport)
	go func() {
		log.Printf("security Serving gRPC on 0.0.0.0:%v", grpcport)
		glog.Fatalln(grpcServer.Serve(lis))
	}()

	//http

	conn, err := grpc.DialContext(context.Background(),
		"0.0.0.0"+":"+grpcport,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Error("Failed to dial server: ", err)
		return
	}

	gwmux := runtime.NewServeMux()
	security.RegisterPkiHandler(context.Background(), gwmux, conn)
	gwServer := &http.Server{
		Addr:    ":" + httpport,
		Handler: gwmux,
	}
	glog.V(2).Info("security Serving gRPC-Gateway on http://0.0.0.0", httpport)
	log.Printf("security Serving gRPC-Gateway on http://0.0.0.0:%v", httpport)
	glog.Fatalln(gwServer.ListenAndServe())
}

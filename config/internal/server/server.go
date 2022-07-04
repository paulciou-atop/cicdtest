package server

import (
	"context"
	"net"
	"net/http"
	"regexp"

	"nms/api/v1/devconfig"
	configServer "nms/config/api/v1/devconfig"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// normalizePort  "134"->":134"| ":8080" -> ":8080" | "abc" -> "abc"
func normalizePort(p string) string {
	var rePort = regexp.MustCompile(`(?m)^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`)
	var rePass = regexp.MustCompile(`(?m)^:([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`)
	if rePass.MatchString(p) {
		return p
	}
	if rePort.MatchString(p) {
		return ":" + p
	}
	return p
}

func RunServer(gRPCPort, httpPort string) error {
	port := normalizePort(gRPCPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {

		log.Error("Failed to listen: ", err)
		return err
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Attach all services to the server
	configServer.RegisterServices(s)
	reflection.Register(s)
	// Serve gRPC Server
	log.Info("Serving gRPC on 0.0.0.0", port)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(context.Background(),
		"0.0.0.0"+port,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to dial server: ", err)
		return err
	}

	gwmux := runtime.NewServeMux()
	port = normalizePort(httpPort)
	err = devconfig.RegisterConfigHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Error("Failed to register gateway:", err)
		return err
	}
	gwServer := &http.Server{
		Addr:    port,
		Handler: gwmux,
	}
	log.Info("Serving gRPC-Gateway on http://0.0.0.0", port)
	log.Fatalln(gwServer.ListenAndServe())

	return nil
}

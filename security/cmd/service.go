package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"security/api/v1/security"
	"security/config"

	"github.com/bobbae/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func NewRunCmd() *cobra.Command {

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "pki service run",
		Run: func(cmd *cobra.Command, args []string) {
			grpcPort, err := cmd.Flags().GetString("grpc-port")
			if err != nil {
				glog.Error("Get grpc-port flag error, use 8080 instead of ")
			}
			httpPort, err := cmd.Flags().GetString("http-port")
			if err != nil {
				glog.Error("Get http-port flag error, use 8090 instead of ")
			}
			lis, err := net.Listen("tcp", ":"+grpcPort)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			grpcServer := grpc.NewServer()

			s := security.NewPkiService()
			security.RegisterPkiServer(grpcServer, s)
			reflection.Register(grpcServer)
			log.Printf("Start service.....%v\n", grpcPort)
			glog.V(2).Infoln("security Serving gRPC on 0.0.0.0", grpcPort)
			go func() {
				log.Printf("security Serving gRPC on 0.0.0.0:%v", grpcPort)
				glog.Fatalln(grpcServer.Serve(lis))
			}()

			//http

			conn, err := grpc.DialContext(context.Background(),
				"0.0.0.0"+":"+grpcPort,
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				glog.Error("Failed to dial server: ", err)
				return
			}

			gwmux := runtime.NewServeMux()
			security.RegisterPkiHandler(context.Background(), gwmux, conn)
			gwServer := &http.Server{
				Addr:    ":" + httpPort,
				Handler: gwmux,
			}
			glog.V(2).Info("security Serving gRPC-Gateway on http://0.0.0.0", httpPort)
			log.Printf("security Serving gRPC-Gateway on http://0.0.0.0:%v", httpPort)
			glog.Fatalln(gwServer.ListenAndServe())
		},
	}
	runCmd.Flags().StringP("grpc-port", "g", config.GetgrpcPort(), "A grpc server listen port")
	runCmd.Flags().StringP("http-port", "p", config.GethttpPort(), "A http server listen port")
	return runCmd
}

package app

import (
	"context"
	"log"
	"net"
	"net/http"
	service "nms/api/v1/atopudpscan"
	"nms/api/v1/configer"
	"nms/atopudpscan/api/v1/atopudpscan"
	"nms/atopudpscan/configs"
	"nms/atopudpscan/internal/pkg/Firmware"

	"github.com/bobbae/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcport = configs.GetgrpcPort()

var httport = configs.GethttpPort()

func NewRunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Service run",
		Run: func(cmd *cobra.Command, args []string) {
			grpcPort, err := cmd.Flags().GetString("grpc-port")
			if err != nil {
				glog.Error(err)

			}
			httpPort, err := cmd.Flags().GetString("http-port")
			if err != nil {
				glog.Error(err)
			}

			lis, err := net.Listen("tcp", ":"+grpcPort)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()
			gwd := atopudpscan.NewGwd()
			err = gwd.Run()
			defer gwd.Close()
			if err != nil {
				glog.Error(err)
			}

			service.RegisterGwdServer(grpcServer, gwd)

			dv := atopudpscan.NewDeviceController()
			service.RegisterAtopDeviceServer(grpcServer, dv)
			conf := atopudpscan.NewConfig()
			configer.RegisterConfigerServer(grpcServer, conf)

			glog.V(2).Infoln("atopudpscan Serving gRPC on 0.0.0.0", grpcPort)
			go func() {
				log.Printf("atopudpscan Serving gRPC on 0.0.0.0:%v", grpcPort)
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

			service.RegisterAtopDeviceHandler(context.Background(), gwmux, conn)
			service.RegisterGwdHandler(context.Background(), gwmux, conn)

			gwServer := &http.Server{
				Addr:    ":" + httpPort,
				Handler: gwmux,
			}
			gwmux.HandlePath("POST", "/api/v1/atopDevice/Upload", Firmware.HandleFileUpload)
			glog.V(2).Info("atopudpscan Serving gRPC-Gateway on http://0.0.0.0", httpPort)
			log.Printf("atopudpscan Serving gRPC-Gateway on http://0.0.0.0:%v", httpPort)
			glog.Fatalln(gwServer.ListenAndServe())
		},
	}

	runCmd.Flags().StringP("grpc-port", "g", grpcport, "A grpc server listen port")
	runCmd.Flags().StringP("http-port", "p", httport, "A http server listen port")
	return runCmd

}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"nms/snmpscan/api/server"

	"github.com/MakeNowJust/heredoc"
	glog "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runHandler(cmd *cobra.Command, args []string) {
	grpcPort, err := cmd.Flags().GetString("grpc-port")
	if err != nil {
		glog.Error("Get grpc-port flag error, use 8080 instead of ")
		grpcPort = "8084"
	}
	httpPort, err := cmd.Flags().GetString("http-port")
	if err != nil {
		glog.Error("Get http-port flag error, use 8090 instead of ")
		httpPort = "8094"
	}
	fmt.Printf("snmpscan service running at gRPC(%s) HTTP(%s)\n", grpcPort, httpPort)
	server.RunGRPCServer(grpcPort, httpPort)
}

// NewRunCmd represents the run command
func NewRunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run snmpscan service",
		Long: heredoc.Doc(`
		Running gRPC & HTTP/HTTPS server`),
		Run: runHandler,
	}
	runCmd.Flags().StringP("grpc-port", "g", "8084", "A grpc server listen port")
	runCmd.Flags().StringP("http-port", "p", "8094", "A http server listen port")
	return runCmd
}

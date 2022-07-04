/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"nms/scanservice/api/server"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run scanservice service",
	Long:  heredoc.Doc("Running gRPC & HTTP/HTTPS server"),
	Run:   runHandler,
}

func runHandler(cmd *cobra.Command, args []string) {
	grpcPort, err1 := cmd.Flags().GetString("grpc-port")
	if err1 != nil {
		grpcPort = "8088"
		glog.Errorf("Get grpc-port flag error, use %v instead of ", grpcPort)
	}
	httpPort, err2 := cmd.Flags().GetString("http-port")
	if err2 != nil {
		httpPort = "8098"
		glog.Errorf("Get http-port flag error, use %v instead of ", httpPort)
	}
	server.RunServer(grpcPort, httpPort)
}

func Run() *cobra.Command {
	// settting grpc port
	runCmd.Flags().StringP("grpc-port", "g", "8088", "A grpc server listen port")
	// setting http port
	runCmd.Flags().StringP("http-port", "p", "8098", "A http server listen port")
	// return run cmd
	return runCmd
}

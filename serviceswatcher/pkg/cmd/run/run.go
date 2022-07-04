/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package run

import (
	"nms/serviceswatcher/api/server"

	"github.com/MakeNowJust/heredoc"
	glog "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runHandler(cmd *cobra.Command, args []string) {
	mode, err := cmd.Flags().GetString("mode")

	if err != nil {
		glog.Error("Get mode flag error, use container instead of ")
		mode = "container"
	}
	viper.Set("mode", mode)
	grpcPort, err := cmd.Flags().GetString("grpc-port")
	if err != nil {
		glog.Error("Get grpc-port flag error, use 8081 instead of ")
		grpcPort = "8081"
	}
	httpPort, err := cmd.Flags().GetString("http-port")
	if err != nil {
		glog.Error("Get http-port flag error, use 8090 instead of ")
		httpPort = "8090"
	}
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
	runCmd.Flags().StringP("grpc-port", "g", "8081", "A grpc server listen port")
	runCmd.Flags().StringP("http-port", "p", "8090", "A http server listen port")
	runCmd.Flags().StringP("mode", "m", "container", "Running mode")

	return runCmd
}

// func init() {
// 	rootCmd.AddCommand(runCmd)

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// runCmd.PersistentFlags().String("grpc-port", "g", "A grpc server listen port")
// 	runCmd.Flags().Int32P("grpc-port", "g", 8080, "A grpc server listen port")
// 	runCmd.Flags().Int32P("http-port", "h", 8090, "A http server listen port")
// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }

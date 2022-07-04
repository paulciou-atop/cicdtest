package run

import (
	"fmt"

	"nms/testing/dummyconfiger/internal/server"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runHandler(cmd *cobra.Command, args []string) {
	grpcPort, err := cmd.Flags().GetString("grpc-port")
	if err != nil {
		log.Error("Get grpc-port flag error, use 8080 instead of ")
		grpcPort = "8085"
	}

	fmt.Printf("config service running at gRPC(%s)\n", grpcPort)
	server.RunServer(grpcPort)
}

func NewCmdRun() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run snmpscan service",
		Long: heredoc.Doc(`
		Running gRPC & HTTP/HTTPS server`),
		Run: runHandler,
	}
	runCmd.Flags().StringP("grpc-port", "g", "8085", "A grpc server listen port")
	return runCmd
}

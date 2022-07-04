package run

import (
	"context"
	"fmt"
	"nms/inventory/internal/server"
	"nms/inventory/pkg/inventory"
	"nms/lib/repo"

	"github.com/MakeNowJust/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runHandler(cmd *cobra.Command, args []string) {
	grpcPort, err := cmd.Flags().GetString("grpc-port")
	if err != nil {
		logrus.Error("Get grpc-port flag error, use 8101 instead of ")
		grpcPort = "8101"
	}

	httpPort, err := cmd.Flags().GetString("http-port")
	if err != nil {
		logrus.Error("Get http-port flag error, use 8111 instead of ")
		httpPort = "8111"
	}

	fmt.Printf("config service running at gRPC(%s) https(%s)\n", grpcPort, httpPort)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := repo.GetRepo(ctx)
	if err != nil {
		logrus.Error("Load repository fail: ", err)
		return
	}
	defer r.Close()

	err = inventory.InitPostgresTable(ctx, r)
	if err != nil {
		logrus.Error(err)
		return
	}

	go inventory.ProcessScanResult(ctx, r)

	server.RunServer(grpcPort, httpPort)

}

func NewCmdRun() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run config service",
		Long: heredoc.Doc(`
		Running gRPC & HTTP/HTTPS server`),
		Run: runHandler,
	}
	runCmd.Flags().StringP("grpc-port", "g", "8101", "A grpc server listen port")
	runCmd.Flags().StringP("http-port", "p", "8111", "A http server listen port")
	return runCmd
}

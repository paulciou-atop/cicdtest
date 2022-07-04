package run

import (
	"fmt"

	"nms/config/internal/server"
	"nms/config/internal/services"
	"nms/config/pkg/config"
	"nms/config/pkg/session"
	"nms/lib/repo"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func runHandler(cmd *cobra.Command, args []string) {
	grpcPort, err := cmd.Flags().GetString("grpc-port")
	if err != nil {
		log.Error("Get grpc-port flag error, use 8100 instead of ")
		grpcPort = "8100"
	}
	httpPort, err := cmd.Flags().GetString("http-port")
	if err != nil {
		log.Error("Get http-port flag error, use 8110 instead of ")
		httpPort = "8110"
	}
	fmt.Printf("config service running at gRPC(%s) HTTP(%s)\n", grpcPort, httpPort)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := repo.GetRepo(ctx)
	services.InitServices(r)
	defer r.Close()
	db := r.DB()
	session.InitDatabaseTables(db)
	config.InitTable(db)
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
	runCmd.Flags().StringP("grpc-port", "g", "8100", "A grpc server listen port")
	runCmd.Flags().StringP("http-port", "p", "8110", "A http server listen port")
	return runCmd
}

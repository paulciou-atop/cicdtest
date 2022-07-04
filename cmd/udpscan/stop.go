package cmd

import (
	"context"

	udpscan "nms/api/v1/atopudpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StopOptions struct {
	ServiceAddr string
	Id          string
}

func NewCmdStop() *cobra.Command {
	o := &StopOptions{}
	// scanCmd represents the scan command
	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "stop udp scan",
		Long: heredoc.Doc(`
		stop udp scan. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				//glog.Fatalf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			c := udpscan.NewGwdClient(conn)

			// Contact the server and print out its response.
			r, err := c.Stop(context.Background(), &udpscan.Sessions{Id: o.Id})
			if err != nil {
				//glog.Fatalf("could not connect service: %v", err)
				return
			}

			cmd.Printf(r.String())
		},
	}

	stopCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	stopCmd.Flags().StringVarP(&o.Id, "session-id", "i", "", "scan id")

	return stopCmd
}

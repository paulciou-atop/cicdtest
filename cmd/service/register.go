package cmd

import (
	"context"
	"log"

	watcher "nms/api/v1/serviceswatcher"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RegsiterOption struct {
	ServiceAddr string
	Name        string //service's name
	Address     string
	Port        int32
	Kind        []string
}

func NewCmdRegister() *cobra.Command {
	o := &RegsiterOption{}
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "register a new service",
		Long: heredoc.Doc(` Atop NMS CLI
		This application is a tool to interface and control the service watcher. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			c := watcher.NewWatcherClient(conn)

			// Contact the server and print out its response.
			r, err := c.Register(context.Background(), &watcher.ServiceInfo{Name: o.Name, Address: o.Address, Port: o.Port, Kind: o.Kind})
			if err != nil {
				log.Fatalf("could not connect watcher service: %v", err)
				return
			}

			cmd.Printf("Register %s %s: %s", o.Name, r.GetSuccess(), r.GetReason())
		},
	}

	registerCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8081", "Service address")
	registerCmd.Flags().StringArrayVarP(&o.Kind, "kind", "k", []string{}, "What kind of API supported")
	registerCmd.Flags().StringVarP(&o.Name, "name", "n", "unknown", "service's name")
	registerCmd.Flags().StringVar(&o.Address, "address", "localhost", "service's hostname/IP address")
	registerCmd.Flags().Int32VarP(&o.Port, "port", "p", 8080, "service's port number")

	registerCmd.MarkFlagRequired("name")
	registerCmd.MarkFlagRequired("address")
	registerCmd.MarkFlagRequired("port")

	return registerCmd
}

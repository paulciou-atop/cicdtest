package cmd

import (
	"context"
	"encoding/json"
	"io"
	"log"

	scanV1 "nms/api/v1/snmpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WalkOptions struct {
	ServiceAddr string
	Target      string
	RootOid     string
}

// NewCmdGet create new get cobra command
func NewCmdWalk() *cobra.Command {
	o := &WalkOptions{}
	var WalkCmd = &cobra.Command{
		Use:   "walk",
		Short: "get subtree",
		Long: heredoc.Doc(`
		Scan online devices which implement SNMP agent. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			c := scanV1.NewSnmpScanClient(conn)

			// Contact the server and print out its response.
			r, err := c.Walkall(context.Background(), &scanV1.WalkallRequest{Target: o.Target, RootOid: o.RootOid})
			if err != nil {
				log.Fatalf("could not connect snmp service: %v", err)
				return
			}

			for {
				stream, err := r.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("error while reading stream: %v\n", err)
					return
				}

				jsonret, err := json.MarshalIndent(stream, "", "  ")
				if err != nil {
					log.Fatalf("could not transfer json format")
					return
				}
				cmd.Printf(string(jsonret))
			}
		},
	}

	WalkCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	WalkCmd.Flags().StringVarP(&o.Target, "target", "t", "127.0.0.1", "Target IP address")
	WalkCmd.Flags().StringVarP(&o.RootOid, "root-oid", "r", ".1.3.6.1.2.1.2.2.1.6", "Root oid")

	viper.BindPFlag("serviceAddr", Cmd.PersistentFlags().Lookup("service-addr"))

	WalkCmd.MarkFlagRequired("target")
	WalkCmd.MarkFlagRequired("root-oid")

	return WalkCmd
}

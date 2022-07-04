package cmd

import (
	"context"
	"encoding/json"
	"log"

	udpscan "nms/api/v1/atopudpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GetServerIpOptions struct {
	ServiceAddr string
	Output      string
}

func NewCmdGetServerIp() *cobra.Command {
	o := &GetServerIpOptions{}
	// scanCmd represents the scan command
	var ipCmd = &cobra.Command{
		Use:   "ip",
		Short: "get server ip",
		Long: heredoc.Doc(`
		get server ip. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(o.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			defer conn.Close()
			c := udpscan.NewGwdClient(conn)

			// Contact the server and print out its response.
			r, err := c.GetServerIp(context.Background(), &udpscan.EmptyParams{})
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			switch {
			case o.Output == "json":
				jsonret, err := json.MarshalIndent(r, "", "  ")
				if err != nil {
					log.Fatalf("could not transfer json format")
					return
				}
				cmd.Printf(string(jsonret))
			default:
				l := list.NewWriter()
				l.SetOutputMirror(cmd.OutOrStderr())
				l.AppendItem("IP")
				l.Indent()
				for _, ip := range r.GetIp() {
					l.AppendItem(ip)
				}
				l.Render()
			}
		},
	}

	ipCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	ipCmd.Flags().StringVarP(&o.Output, "output format", "o", "", "Output result with format, support json")

	return ipCmd
}

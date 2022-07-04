package cmd

import (
	"context"
	"encoding/json"
	"log"

	udpscan "nms/api/v1/atopudpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SessionDataOptions struct {
	ServiceAddr string
	Id          string
	Format      string
}

func NewCmdSessionData() *cobra.Command {
	o := &SessionDataOptions{}
	// sessionDataCmd represents the scan command
	var sessionDataCmd = &cobra.Command{
		Use:   "sessiondata",
		Short: "gwd scan result",
		Long: heredoc.Doc(`
		Gwd scan result. 
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
			r, err := c.GetSessionData(context.Background(), &udpscan.Sessions{Id: o.Id})
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			switch {
			case o.Format == "json":
				jsonret, err := json.MarshalIndent(r, "", "  ")
				if err != nil {
					log.Fatalf("could not transfer json format")
					return
				}
				cmd.Printf(string(jsonret))
			default:
				t := table.NewWriter()
				t.SetStyle(table.StyleRounded)
				t.SetOutputMirror(cmd.OutOrStderr())
				t.AppendHeader(table.Row{"Model", "MAC", "IP", "Net Mask", "Gateway", "Host", "Kernel", "AP", "DHCP"})
				for _, dev := range r.GetDevices() {
					t.AppendRow(table.Row{dev.GetModel(), dev.GetMacAddress(), dev.GetIPAddress(), dev.GetNetmask(), dev.GetGateway(), dev.GetHostname(), dev.GetKernel(), dev.GetAp(), dev.GetIsDHCP()})
				}
				t.Render()
			}
		},
	}

	sessionDataCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	sessionDataCmd.Flags().StringVarP(&o.Id, "session-id", "i", "", "scan id")
	sessionDataCmd.Flags().StringVarP(&o.Format, "output format", "f", "", "Output result with format, support json")

	viper.BindPFlag("gwd.server", Cmd.PersistentFlags().Lookup("server-ip"))

	sessionDataCmd.MarkFlagRequired("server-ip")
	return sessionDataCmd
}

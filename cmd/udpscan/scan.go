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

type ScanOptions struct {
	ServiceAddr string
	ServerIp    string
	Id          string
	Format      string
}

func NewCmdScan() *cobra.Command {
	o := &ScanOptions{}
	// scanCmd represents the scan command
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scan atop devices",
		Long: heredoc.Doc(`
		Scan online atop devices. 
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
			r, err := c.SessionScan(context.Background(), &udpscan.ScanConfig{ServerIp: o.ServerIp, Id: o.Id})
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
				t.AppendHeader(table.Row{"ID", "Status", "Message"})
				t.AppendRow(table.Row{r.GetId(), r.GetStatus(), r.GetMessage()})
				t.Render()
			}

		},
	}

	scanCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	scanCmd.Flags().StringVarP(&o.ServerIp, "server-ip", "s", "", "server receive ip")
	scanCmd.Flags().StringVarP(&o.Id, "session-id", "i", "", "scan id")
	scanCmd.Flags().StringVarP(&o.Format, "output format", "f", "", "Output result with format, support json")

	viper.BindPFlag("gwd.server", Cmd.PersistentFlags().Lookup("server-ip"))

	scanCmd.MarkFlagRequired("server-ip")
	return scanCmd
}

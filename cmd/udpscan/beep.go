package cmd

import (
	"context"

	udpscan "nms/api/v1/atopudpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BeepOptions struct {
	ServiceAddr string
	DeviceIp    string
	MACAddress  string
	ServerIp    string
}

func NewCmdBeep() *cobra.Command {
	o := &BeepOptions{}
	// scanCmd represents the scan command
	var beepCmd = &cobra.Command{
		Use:   "beep [deviceip] [devicemac] [severip]",
		Short: "make device sound",
		Long: heredoc.Doc(`
		make device sound. 
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
			r, err := c.Beep(context.Background(), &udpscan.GwdConfig{IPAddress: o.DeviceIp, MACAddress: o.MACAddress, ServerIp: o.ServerIp})
			if err != nil {
				//glog.Fatalf("could not connect service: %v", err)
				return
			}

			cmd.Printf(r.String())
		},
	}

	beepCmd.Flags().StringVar(&o.ServiceAddr, "service-addr", "127.0.0.1:8080", "Service address")
	beepCmd.Flags().StringVarP(&o.DeviceIp, "device", "d", "", "device ip")
	beepCmd.Flags().StringVarP(&o.MACAddress, "mac", "m", "", "device mac")
	beepCmd.Flags().StringVarP(&o.ServerIp, "server-ip", "s", "", "server ip")

	viper.BindPFlag("gwd.deviceip", Cmd.PersistentFlags().Lookup("device"))
	viper.BindPFlag("gwd.macaddress", Cmd.PersistentFlags().Lookup("mac"))
	viper.BindPFlag("gwd.server", Cmd.PersistentFlags().Lookup("server"))

	beepCmd.MarkFlagRequired("server-ip")
	return beepCmd
}

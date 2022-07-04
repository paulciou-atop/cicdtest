package app

import (
	"context"
	"log"
	"net"
	"nms/api/v1/atopudpscan"
	"nms/atopudpscan/configs"
	atopnet "nms/atopudpscan/pkg/net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	NewBeepCmd()
	NewSessionScanCmd()
	NewGetServerIpCmd()
	NewstopCmd()
}

var GwdCmd = &cobra.Command{
	Use:   "gwd",
	Short: "Gwd function",
	Long:  "Gwd function",
}

func NewBeepCmd() {
	beepCmd := &cobra.Command{
		Use:   "beep",
		Short: "make device sound",
		Long:  "make device sound",
		Run: func(cmd *cobra.Command, args []string) {

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			device, err := cmd.Flags().GetString("device")
			if err != nil {
				log.Fatalln(err)

			}

			mac, err := cmd.Flags().GetString("mac")
			if err != nil {
				log.Fatalln(err)

			}
			out, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server," + err.Error())

			}
			err = atopnet.CheckIPAddress(device)
			if err != nil {
				log.Fatalf("device," + err.Error())

			}
			err = atopnet.CheckMacAddress(mac)
			if err != nil {
				log.Fatalf("mac," + err.Error())

			}
			err = atopnet.CheckIPAddress(out)
			if err != nil {
				log.Fatalf("out," + err.Error())

			}
			beep(server, port, device, mac, out)

		},
	}
	beepCmd.Flags().StringP("server", "s", "", "server ip")
	beepCmd.Flags().StringP("port", "p", grpcport, "server port")
	beepCmd.Flags().StringP("device", "d", "", "device ip")
	beepCmd.Flags().StringP("mac", "m", "", "device mac")
	beepCmd.Flags().StringP("out", "o", "", "out ip path of server")
	GwdCmd.AddCommand(beepCmd)
}

func beep(serverip, port, deviceip, macaddress, out string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(serverip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.GwdConfig{IPAddress: deviceip, MACAddress: macaddress, ServerIp: out}
	r, err := u.Beep(ctx, n)
	if err != nil {
		log.Fatalf("could not Beep: %v", err)

	}
	log.Println(r)
}

func NewGetServerIpCmd() {
	var serverIpCmd = &cobra.Command{
		Use:   "ip",
		Short: "get server ip",
		Long:  "get server ip",
		Run: func(cmd *cobra.Command, args []string) {

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server:" + err.Error())

			}
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			addr := net.JoinHostPort(server, port)
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			u := atopudpscan.NewGwdClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			r, err := u.GetServerIp(ctx, &atopudpscan.EmptyParams{})
			if err != nil {
				log.Fatalf("could not GetServerIp: %v", err)

			}
			for _, v := range r.GetIp() {
				log.Println(v)
			}
		},
	}
	serverIpCmd.Flags().StringP("server", "s", "", "server ip")
	serverIpCmd.Flags().StringP("port", "p", grpcport, "server port")
	GwdCmd.AddCommand(serverIpCmd)
}

func NewSessionScanCmd() {
	sessionscanCmd := &cobra.Command{
		Use:   "sessionscan",
		Short: "session scan device",
		Long:  "contronl service to scan device in session",
		Run: func(cmd *cobra.Command, args []string) {
			serverip, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			out, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatalln(err)

			}
			id, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatalln(err)

			}

			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			sessionscan(serverip, port, out, id)
		},
	}

	sessionscanCmd.Flags().StringP("server", "s", "", "server ip")
	sessionscanCmd.Flags().StringP("out", "o", "", "out ip path of server")
	sessionscanCmd.Flags().StringP("id", "i", "1", "session id")
	sessionscanCmd.Flags().StringP("port", "p", configs.GetgrpcPort(), "server port")

	GwdCmd.AddCommand(sessionscanCmd)

}

func sessionscan(serverip string, port string, out string, id string) {
	err := atopnet.CheckIPAddress(serverip)
	if err != nil {
		log.Fatal("serverIp," + err.Error())

	}
	err = atopnet.CheckIPAddress(out)
	if err != nil {
		log.Fatal("out", err.Error())

	}
	addr := net.JoinHostPort(serverip, port)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	n := &atopudpscan.ScanConfig{ServerIp: out, Id: id}
	r, err := u.SessionScan(ctx, n)
	if err != nil {
		log.Fatalf("could not SessionScan: %v", err)
	}
	log.Printf("%v", r)

}

func NewstopCmd() {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "stop scan of device",
		Long:  "stop scan of device",
		Run: func(cmd *cobra.Command, args []string) {
			serverip, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}

			id, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatalln(err)

			}

			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			addr := net.JoinHostPort(serverip, port)
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			u := atopudpscan.NewGwdClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			r, err := u.Stop(ctx, &atopudpscan.Sessions{Id: id})
			if err != nil {
				log.Fatalf("could not Scan: %v", err)

			}
			log.Println(r)

		},
	}

	stopCmd.Flags().StringP("server", "s", "", "server ip")
	stopCmd.Flags().StringP("id", "i", "1", "session id")
	stopCmd.Flags().StringP("port", "p", configs.GetgrpcPort(), "server port")

	GwdCmd.AddCommand(stopCmd)
}

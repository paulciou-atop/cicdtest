package app

import (
	"context"
	"io"
	"log"
	"net"
	"nms/api/v1/atopudpscan"
	"nms/atopudpscan/configs"
	"os"
	"time"

	atopnet "nms/atopudpscan/pkg/net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "device function",
	Long:  "device function",
}

func init() {
	NewRebootCmd()

	NewSettingConfigCmd()
	NewfwupgradingCmd()
	NewUploadCmd()
	NewfwupStatusCmd()
}
func NewRebootCmd() {
	var rebootCmd = &cobra.Command{
		Use:   "reboot",
		Short: "reboot device",
		Long:  "reboot device",
		Run: func(cmd *cobra.Command, args []string) {

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			out, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatalln(err)

			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			mac, err := cmd.Flags().GetString("mac")
			if err != nil {
				log.Fatalln(err)

			}
			deivce, err := cmd.Flags().GetString("deivce")
			if err != nil {
				log.Fatalln(err)

			}

			user, err := cmd.Flags().GetString("user")
			if err != nil {
				log.Fatalln(err)

			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				log.Fatalln(err)

			}

			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server," + err.Error())
			}

			err = atopnet.CheckIPAddress(out)
			if err != nil {
				log.Fatalf("out," + err.Error())
			}

			err = atopnet.CheckMacAddress(mac)
			if err != nil {
				log.Fatalf("mac," + err.Error())
			}

			err = atopnet.CheckIPAddress(deivce)
			if err != nil {
				log.Fatalf("deivce," + err.Error())
			}

			reboot(server, out, port, mac, deivce, user, password)
		},
	}
	rebootCmd.Flags().StringP("server", "s", "", "server ip")
	rebootCmd.Flags().StringP("out", "o", "", "out ip path of server")
	rebootCmd.Flags().StringP("port", "p", grpcport, "server port")
	rebootCmd.Flags().StringP("mac", "m", "", "device mac")
	rebootCmd.Flags().StringP("deivce", "d", "", "device ip")
	rebootCmd.Flags().StringP("user", "u", account, "device user name")
	rebootCmd.Flags().StringP("password", "P", password, "device user password")
	deviceCmd.AddCommand(rebootCmd)
}

func reboot(serverip, scannet, port, mac, deviceIp, userName, passWord string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(serverip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.GwdConfig{IPAddress: deviceIp,
		MACAddress: mac,
		ServerIp:   scannet,
		Username:   userName,
		Password:   passWord}
	_, err = u.Reboot(ctx, n)

	if err != nil {
		log.Fatalf("could not Reboot: %v", err)

	}
}

func NewSettingConfigCmd() {

	var settingCmd = &cobra.Command{
		Use:   "setting",
		Short: "setting device config",
		Long:  "setting device config",
		Run: func(cmd *cobra.Command, args []string) {

			device, err := cmd.Flags().GetString("device")
			if err != nil {
				log.Fatalln(err)

			}
			mac, err := cmd.Flags().GetString("mac")
			if err != nil {
				log.Fatalln(err)

			}
			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			user, err := cmd.Flags().GetString("user")
			if err != nil {
				log.Fatalln(err)

			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				log.Fatalln(err)

			}
			newip, err := cmd.Flags().GetString("newip")
			if err != nil {
				log.Fatalln(err)

			}
			mask, err := cmd.Flags().GetString("mask")
			if err != nil {
				log.Fatalln(err)

			}
			gateway, err := cmd.Flags().GetString("gateway")
			if err != nil {
				log.Fatalln(err)

			}
			hostname, err := cmd.Flags().GetString("hostname")
			if err != nil {
				log.Fatalln(err)

			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			out, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(device)
			if err != nil {
				log.Fatalf("device," + err.Error())
			}
			err = atopnet.CheckMacAddress(mac)
			if err != nil {
				log.Fatalf("mac," + err.Error())
			}
			err = atopnet.CheckIPAddress(newip)
			if err != nil {
				log.Fatalf("newip," + err.Error())
			}
			err = atopnet.CheckIPAddress(mask)
			if err != nil {
				log.Fatalf("mask" + err.Error())
			}
			err = atopnet.CheckIPAddress(gateway)
			if err != nil {
				log.Fatalf("gateway," + err.Error())
			}
			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server," + err.Error())
			}

			err = atopnet.CheckIPAddress(out)
			if err != nil {
				log.Fatalf("out," + err.Error())
			}
			n := &atopudpscan.GwdConfig{
				IPAddress:    device,
				MACAddress:   mac,
				NewIPAddress: newip,
				Netmask:      mask,
				Gateway:      gateway,
				Hostname:     hostname,
				Username:     user,
				Password:     password,
				ServerIp:     out,
			}
			setting(n, server, port)
		},
	}
	settingCmd.Flags().StringP("device", "d", "", "device ip")
	settingCmd.Flags().StringP("mac", "m", "", "device macaddress")
	settingCmd.Flags().StringP("server", "s", "", "server ip")
	settingCmd.Flags().StringP("user", "u", account, "device user name")
	settingCmd.Flags().StringP("password", "P", password, "device user password")
	settingCmd.Flags().StringP("newip", "n", "", "new ip of device new ip")
	settingCmd.Flags().StringP("mask", "M", "", "new netmask of device ")
	settingCmd.Flags().StringP("gateway", "g", "", "new gateway of device")
	settingCmd.Flags().StringP("hostname", "H", "", "device hostname")
	settingCmd.Flags().StringP("port", "p", grpcport, "port")
	settingCmd.Flags().StringP("out", "o", "", "out ip path of server")
	deviceCmd.AddCommand(settingCmd)
}

func setting(a *atopudpscan.GwdConfig, serverip, port string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(serverip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)

	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r, err := u.SettingConfig(ctx, a)
	if err != nil {
		log.Fatalf("could not SettingConfig: %v", err)

	}
	log.Print(r)
}

func NewfwupgradingCmd() {
	var fwUpgradingCmd = &cobra.Command{
		Use:   "fwupgrade",
		Short: "device updrade fw",
		Long:  "device updrade fw after complete file deleted",
		Run: func(cmd *cobra.Command, args []string) {
			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			device, err := cmd.Flags().GetString("device")
			if err != nil {
				log.Fatalln(err)

			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			filename, err := cmd.Flags().GetString("filename")
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
			n := &atopudpscan.FwInfo{
				DeviceIp: device,
				FileName: filename,
			}
			fwupGrading(n, server, port)
		},
	}
	fwUpgradingCmd.Flags().StringP("server", "s", "", "server ip")
	fwUpgradingCmd.Flags().StringP("device", "d", "", "device ip")
	fwUpgradingCmd.Flags().StringP("port", "p", grpcport, "port")
	fwUpgradingCmd.Flags().StringP("filename", "f", "", "filename")
	deviceCmd.AddCommand(fwUpgradingCmd)

}

func fwupGrading(a *atopudpscan.FwInfo, serverip, port string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(serverip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)

	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	r, err := u.FwUpgrading(ctx, a)
	if err != nil {
		log.Fatalf("could not FwUpgrading: %v", err)

	}
	log.Print(r)
}

func GetProcessStatus(device string, serverip, port string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(serverip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.FwRequest{
		IPAddress: device,
	}
	r, err := u.GetProcessStatus(ctx, n)
	if err != nil {
		log.Fatalf("could not GetProcessStatus: %v\n", err)

	}
	log.Println(r)
}

func NewfwupStatusCmd() {

	var fwStatusCmd = &cobra.Command{
		Use:   "fwstatus",
		Short: "get device fw upgrade status",
		Long:  "get device fw upgrade status",
		Run: func(cmd *cobra.Command, args []string) {
			device, err := cmd.Flags().GetString("device")
			if err != nil {
				log.Fatalln(err)

			}

			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(device)
			if err != nil {
				log.Fatalf("device," + err.Error())
			}

			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server," + err.Error())
			}

			GetProcessStatus(device, server, port)

		},
	}
	fwStatusCmd.Flags().StringP("device", "d", "", "device ip")
	fwStatusCmd.Flags().StringP("server", "s", "", "server ip")
	fwStatusCmd.Flags().StringP("port", "p", configs.GetgrpcPort(), "port")
	deviceCmd.AddCommand(fwStatusCmd)
}

func NewUploadCmd() {
	UploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload file",
		Long:  "Upload file to service",
		Run: func(cmd *cobra.Command, args []string) {
			serverip, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				log.Fatalln(err)

			}

			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			upload(serverip, port, file)
		},
	}

	UploadCmd.Flags().StringP("server", "s", "", "server ip")
	UploadCmd.Flags().StringP("file", "f", "", "file name")
	UploadCmd.Flags().StringP("port", "p", configs.GetgrpcPort(), "server port")

	deviceCmd.AddCommand(UploadCmd)
}

func upload(serverip, port, filename string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(serverip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)

	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not FwUpgrading: %v", err)
	}
	stream, err := u.Upload(ctx)
	if err != nil {
		log.Fatalf("could not FwUpgrading: %v", err)
	}
	buf := make([]byte, 1024)

	for {
		num, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("could not FwUpgrading: %v", err)
		}
		if err := stream.Send(&atopudpscan.Chunk{Content: buf[:num]}); err != nil {
			log.Printf("could not FwUpgrading: %v", err)
		}

	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)

	}
	log.Printf("result :%v\n", resp.Code)
	log.Printf("message :%v\n", resp.Message)
	log.Printf("filename :%v\n", resp.FileName)
}

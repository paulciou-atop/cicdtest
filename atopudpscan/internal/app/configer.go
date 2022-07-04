package app

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"nms/api/v1/common"
	"nms/api/v1/configer"
	"nms/api/v1/devconfig"
	atopnet "nms/atopudpscan/pkg/net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	ConfigCmd.AddCommand(NewSetting())
	ConfigCmd.AddCommand(NewGetSetting())
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "config function",
	Long:  "config function",
}

//impelmenting.....
func NewSetting() *cobra.Command {
	var set = &cobra.Command{
		Use:   "setting",
		Short: "configure  setting",
		Run: func(cmd *cobra.Command, args []string) {
			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server," + err.Error())
			}

			out, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(out)
			if err != nil {
				log.Fatalf("out ip," + err.Error())
			}
			mac, err := cmd.Flags().GetString("mac")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckMacAddress(mac)
			if err != nil {
				log.Fatalf("mac," + err.Error())
			}
			newip, err := cmd.Flags().GetString("newip")
			if err != nil {
				log.Fatalln(err)

			}
			if newip != "" {
				err = atopnet.CheckIPAddress(newip)
				if err != nil {
					log.Fatalf("New ip," + err.Error())
				}
			}

			mask, err := cmd.Flags().GetString("mask")
			if err != nil {
				log.Fatalln(err)

			}
			if mask != "" {
				err = atopnet.CheckIPAddress(mask)
				if err != nil {
					log.Fatalf("mask," + err.Error())
				}
			}

			gateway, err := cmd.Flags().GetString("gateway")
			if err != nil {
				log.Fatalln(err)

			}
			if gateway != "" {
				err = atopnet.CheckIPAddress(gateway)
				if err != nil {
					log.Fatalf("gateway," + err.Error())
				}
			}

			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}

			id, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatalln(err)

			}
			hostname, err := cmd.Flags().GetString("hostname")
			if err != nil {
				log.Fatalln(err)

			}
			m := make(map[string]interface{}, 0)
			m["Username"] = account
			m["Password"] = password
			m["ServerIp"] = out
			m["Hostname"] = hostname
			m["NewIPAddress"] = newip
			m["Netmask"] = mask
			m["Gateway"] = gateway
			config(server, port, m, mac, id)
		},
	}

	set.Flags().StringP("server", "s", "", "server ip")
	set.Flags().StringP("port", "p", grpcport, "server port")
	set.Flags().StringP("out", "o", "0.0.0.0", "out ip path of server")
	set.Flags().StringP("mac", "m", "", "device macaddress")
	set.Flags().StringP("user", "u", account, "device user name")
	set.Flags().StringP("password", "P", password, "device user password")
	set.Flags().StringP("newip", "n", "", "new ip of device new ip")
	set.Flags().StringP("mask", "M", "", "new netmask of device ")
	set.Flags().StringP("gateway", "g", "", "new gateway of device")
	set.Flags().StringP("hostname", "H", "", "device hostname")
	set.Flags().StringP("id", "i", "1", "session id")
	return set

}

func config(ip, port string, m map[string]interface{}, mac string, id string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	addr := net.JoinHostPort(ip, port)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := configer.NewConfigerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	b, err := json.Marshal(m)
	sess := &devconfig.SessionState{Id: id}
	p := &structpb.Struct{}
	err = protojson.Unmarshal(b, p)
	conf := &devconfig.ConfigOptions{Protocol: "gwd", Kind: "network", Payload: p}
	confs := make([]*devconfig.ConfigOptions, 0)
	confs = append(confs, conf)
	d := &common.DeviceIdentify{DeviceId: mac}
	n := &configer.ConfigerConfigRequest{
		Session: sess,
		Configs: confs,
		Device:  d,
	}
	stream, err := u.Config(ctx)
	if err == io.EOF {
		return
	}
	err = stream.Send(n)
	if err != nil {
		log.Fatal(err)
		return
	}
	for {

		r, err := stream.Recv()
		if err == io.EOF {

			return
		}
		if err != nil {
			log.Print(err)
			return
		}
		if r.Session.State == "running" {
			log.Print(r)
		}

		if r.Session.State == "success" {
			log.Print(r)
			stream.CloseSend()

		} else if r.Session.State == "fail" {
			log.Print(r)
			stream.CloseSend()
		}

	}

}

func NewGetSetting() *cobra.Command {
	var get = &cobra.Command{
		Use:   "get",
		Short: "Get setting",
		Run: func(cmd *cobra.Command, args []string) {
			server, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckIPAddress(server)
			if err != nil {
				log.Fatalf("server," + err.Error())
			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalln(err)

			}
			mac, err := cmd.Flags().GetString("mac")
			if err != nil {
				log.Fatalln(err)

			}
			err = atopnet.CheckMacAddress(mac)
			if err != nil {
				log.Fatalf("mac," + err.Error())
			}

			addr := net.JoinHostPort(server, port)
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			u := configer.NewConfigerClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			d := &common.DeviceIdentify{DeviceId: mac}
			n := &configer.GetConfigRequest{
				Device: d,
			}
			r, err := u.GetConfig(ctx, n)
			if err != nil {
				log.Fatal(err)
			}
			b, _ := protojson.Marshal(r)
			log.Print(b)
		},
	}

	get.Flags().StringP("server", "s", "", "server ip")
	get.Flags().StringP("port", "p", grpcport, "server port")
	get.Flags().StringP("mac", "m", "", "device macaddress")

	return get

}

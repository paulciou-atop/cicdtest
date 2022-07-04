package main_test

import (
	"context"
	"io"
	"log"
	"nms/api/v1/atopudpscan"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestReboot(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.GwdConfig{IPAddress: "192.168.4.30",
		MACAddress: "00-60-e9-18-3c-3c",
		ServerIp:   "192.168.4.21",
		Username:   "admin",
		Password:   "default"}
	_, err = u.Reboot(ctx, n)

	if err != nil {
		log.Fatalf("could not Reboot: %v", err)

	}
}

func TestSettingConfig(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	n := &atopudpscan.GwdConfig{
		IPAddress:    "192.168.4.30",
		MACAddress:   "00-60-e9-18-3c-3c",
		NewIPAddress: "192.168.4.30",
		Netmask:      "255.255.255.0",
		Gateway:      "192.168.4.254",
		Hostname:     "atop13",
		Username:     "admin",
		Password:     "default",
		ServerIp:     "192.168.4.21"}
	_, err = u.SettingConfig(ctx, n)
	if err != nil {
		log.Fatalf("could not SettingConfig: %v", err)

	}

}

func TestFwUpgrading(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	n := &atopudpscan.FwInfo{
		DeviceIp: "192.168.4.30",
		FileName: "20220518132529.dld",
	}
	r, err := u.FwUpgrading(ctx, n)

	if err != nil {
		log.Fatalf("could not FwUpgrading: %v\n", err)

	}
	log.Printf("Result:%v", r.Result)
	log.Printf("Message:%v", r.Message)

	/*
		n := &atopudpscan.FwConfig{
			IPAddress: "192.168.4.30",
			FileName:  "EH7520-4G-8PoE-4SFP_K544_A544.dld",
		}
		go func() {
			for {
				time.Sleep(time.Second * 1)
				GetProcessStatus()
			}
		}()
		_, err = u.FwUpgrading(ctx, n)
		if err != nil {
			log.Fatalf("could not FwUpgrading: %v", err)

		}*/

}

func TestUpload(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()
	f, err := os.Open("EH7520-4G-8PoE-4SFP_K545_A545.dld")
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

func GetProcessStatus() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.FwRequest{
		IPAddress: "192.168.4.30",
	}
	r, err := u.GetProcessStatus(ctx, n)
	if err != nil {
		log.Fatalf("could not GetProcessStatus: %v\n", err)
	}
	log.Println(r.GetMessage())

}

func TestDeviceServerIp(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	r, err := u.GetServerIp(ctx, &atopudpscan.EmptyParams{})
	if err != nil {
		log.Fatalf("could not DeviceServerIp: %v", err)

	}
	for _, v := range r.GetIp() {
		log.Println(v)
	}
}

/*
func TestDleteFile(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewAtopDeviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.FwInfo{
		FileName: "20220523135536.dld",
	}
	r, err := u.DeleteFile(ctx, n)
	if err != nil {
		log.Fatalf("could not DleteFile: %v", err)

	}
	log.Println(r)
}
*/

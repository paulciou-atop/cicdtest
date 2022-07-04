package main_test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"nms/api/v1/common"
	"nms/api/v1/configer"
	"nms/atopudpscan/pkg/Simulatedevice"
	"strconv"
	"sync"
	"testing"
	"time"

	"nms/api/v1/devconfig"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestVailed(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//u := atopudpscan.NewConfigurationClient(conn)
	u := configer.NewConfigerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	m := map[string]interface{}{
		"MACAddress": "bar",
		"Hostname":   123,
	}

	b, err := json.Marshal(m)
	sess := &devconfig.SessionState{Id: "1"}
	p := &structpb.Struct{}
	err = protojson.Unmarshal(b, p)
	conf := &devconfig.ConfigOptions{Protocol: "gwd", Kind: "network", Payload: p}
	confs := make([]*devconfig.ConfigOptions, 0)
	confs = append(confs, conf)
	n := &configer.ConfigerValidateRequest{
		Session: sess,
		Configs: confs,
	}
	stream, err := u.Validate(ctx)

	if err == io.EOF {
		log.Print(err)
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
		if r.Session.State == "success" {
			log.Print(r)
			stream.CloseSend()

		} else if r.Session.State == "fail" {
			log.Print(r)
			stream.CloseSend()
		}

	}

}

func TestConfig(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := configer.NewConfigerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	m := map[string]interface{}{
		"NewIPAddress": "192.168.4.51",
		"Username":     "admin",
		"Hostname":     "sw78988",
		"Password":     "default",
		"ServerIp":     "0.0.0.0",
	}

	b, err := json.Marshal(m)
	sess := &devconfig.SessionState{Id: "1"}
	p := &structpb.Struct{}
	err = protojson.Unmarshal(b, p)
	conf := &devconfig.ConfigOptions{Protocol: "gwd", Kind: "network", Payload: p}
	confs := make([]*devconfig.ConfigOptions, 0)
	confs = append(confs, conf)
	d := &common.DeviceIdentify{DeviceId: "00-60-E9-2C-01-01"}
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

const number = 29

func TestMultipleConfig(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := configer.NewConfigerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	v := getTestmutipleConfig(number)

	for _, v := range v {
		ch := make(chan bool)
		stream, err := u.Config(ctx)
		if err == io.EOF {
			return
		}
		err = stream.Send(&v)
		if err != nil {
			log.Fatal(err)
			return
		}
		go func() {
			wg.Add(1)
			defer wg.Done()
			ch <- true
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

		}()
		<-ch
	}
	wg.Wait()
}

func getTestmutipleConfig(number int) []configer.ConfigerConfigRequest {

	res := make([]configer.ConfigerConfigRequest, 0)
	mp := make([]map[string]interface{}, 0)
	for i := 1; i <= number; i++ {
		d, _ := Simulatedevice.GetTestParam("device", uint(i))
		n := Simulatedevice.GetTestHostName("tt", i)
		ip := Simulatedevice.GetTestNewIP(i)
		g := Simulatedevice.GetGateway(i)
		msk := Simulatedevice.GetTestNetMask(i)
		m := map[string]interface{}{
			"NewIPAddress": ip,
			"Netmask":      msk,
			"Gateway":      g,
			"Username":     "admin",
			"Hostname":     n,
			"Password":     "default",
			"ServerIp":     "0.0.0.0",
			"MACAddress":   d.MACAddress,
		}
		mp = append(mp, m)
	}
	for k, v := range mp {
		b, _ := json.Marshal(v)
		sess := &devconfig.SessionState{Id: strconv.Itoa(k)}
		p := &structpb.Struct{}
		protojson.Unmarshal(b, p)
		conf := &devconfig.ConfigOptions{Protocol: "gwd", Kind: "network", Payload: p}
		confs := make([]*devconfig.ConfigOptions, 0)
		confs = append(confs, conf)
		confs = append(confs, conf)
		d := &common.DeviceIdentify{DeviceId: v["MACAddress"].(string)}
		n := configer.ConfigerConfigRequest{
			Session: sess,
			Configs: confs,
			Device:  d,
		}
		res = append(res, n)
	}
	return res
}

/*
func TestFileTransfer(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewConfigurationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	sess := &atopudpscan.SessionState{
		Id: "00-60-e9-18-3c-3c",
	}
	f := &atopudpscan.FileParam{FileLocation: "https://github.com/pcman-bbs/pcman-windows/releases/download/v9.5.0-beta3/PCMan.exe", DeviceId: "192.168.4.30", Upload: true}

	n := &atopudpscan.ClientFileTransferRequest{
		Session:   sess,
		FileParam: f,
	}

	stream, err := u.FileTransfer(ctx)

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
*/
func TestGetFile(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := configer.NewConfigerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	d := &common.DeviceIdentify{DeviceId: "00-60-e9-20-0a-44"}
	n := &configer.GetConfigRequest{
		Device: d,
	}
	r, err := u.GetConfig(ctx, n)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(r)
}

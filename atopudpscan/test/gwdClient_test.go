package main_test

import (
	"context"
	"log"
	"nms/api/v1/atopudpscan"
	"nms/atopudpscan/configs"
	"testing"
	"time"

	"github.com/Atop-NMS-team/pgutils"

	"google.golang.org/grpc"
)

var (
	address = "localhost:" + configs.GetgrpcPort()
)

/*
func TestScan(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5*timeout)) //suggest:time out should more than 5 seconde of scan timeout,beacuse scan include savedatebase
	defer cancel()
	n := &atopudpscan.ScanConfig{ServerIp: "192.168.4.21", Id: "2"}
	r, err := u.Scan(ctx, n)
	if err != nil {
		log.Fatalf("could not Scan: %v", err)

	}
	for _, v := range r.GetDevices() {
		log.Println(v)

	}
	//log.Println(string(r.GetDevices()))
}
*/
func TestGwdStop(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	n := &atopudpscan.ScanConfig{ServerIp: "192.168.4.244", Id: "1"}
	go func() {
		time.Sleep(time.Second * 1)
		r, err := u.Stop(ctx, &atopudpscan.Sessions{Id: "1"})
		if err != nil {
			log.Fatalf("could not Scan: %v", err)

		}
		log.Println(r)
	}()
	r, err := u.SessionScan(ctx, n)
	if err != nil {
		log.Fatalf("could not Scan: %v", err)

	}
	log.Println(r)
}

func TestBeep(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	n := &atopudpscan.GwdConfig{IPAddress: "192.168.4.30", MACAddress: "00-60-e9-18-3c-3c", ServerIp: "192.168.4.21"}
	r, err := u.Beep(ctx, n)
	if err != nil {
		log.Fatalf("could not Beep: %v", err)

	}
	log.Println(r)
}

func TestGwdServerIp(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
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
}

func TestScanSession(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30) //suggest:time out should more than 5 seconde of scan timeout,beacuse scan include savedatebase
	defer cancel()
	c, err := pgutils.NewClient()
	if err != nil {
		log.Print(err)

	} else {
		c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})
		Insert("gwd:running", "1")
	}
	n := &atopudpscan.ScanConfig{ServerIp: "0.0.0.0", Id: "1"}
	r, err := u.SessionScan(ctx, n)
	if err != nil {
		log.Fatalf("could not SessionScan: %v", err)

	}
	if err != nil {
		log.Fatalf("could not SessionScan: %v", err)
	}
	log.Printf("%v", r)
}

func TestGetSessionStatus(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	r, err := u.GetSessionStatus(ctx, &atopudpscan.Sessions{Id: "2"})
	if err != nil {
		log.Fatalf("could GetSessionStatus: %v", err)

	}

	log.Printf("%v", r)
}

func TestGetessionDate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	u := atopudpscan.NewGwdClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	r, err := u.GetSessionData(ctx, &atopudpscan.Sessions{Id: "1"})
	if err != nil {
		log.Fatalf("could GetessionDate: %v", err)

	}
	log.Printf("%v", r)

}

func TestSql(t *testing.T) {
	Insert("gwd:running", "1")
}

func Insert(gwd string, id string) {
	v := gwd + "|snmp: fail"
	c, err := pgutils.NewClient()
	if err != nil {
		log.Print(err)
		return
	}
	c.CreateTable(&pgutils.DeviceSession{}, pgutils.CreateTableOpt{IfNotExists: true})
	c.Insert(&pgutils.DeviceSession{
		SessionID:   id,
		State:       v,
		CreatedTime: time.Now().String(),
	})
	if err != nil {
		log.Print(err)
	}
}

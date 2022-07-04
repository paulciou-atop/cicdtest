package main

import (
	"flag"
	"fmt"
	"net"

	scanV1 "nms/api/v1/snmpscan"

	"google.golang.org/grpc"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

var (
	port = flag.Int("port", 40051, "The server test port")
)

type Server struct {
	scanV1.UnimplementedSnmpScanServer
}

func (s *Server) Get(_ *scanV1.GetRequest, stream scanV1.SnmpScan_GetServer) error {
	fmt.Println("get request")
	pbV, err := structpb.NewValue("get test")
	if err != nil {
		return err
	}

	pdu := scanV1.PDU{Value: pbV, Name: "Test"}

	if err = stream.Send(&pdu); err != nil {
		return err
	}

	fmt.Println("Sending test: get test")
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	s := grpc.NewServer()
	scanV1.RegisterSnmpScanServer(s, &Server{})
	fmt.Println("server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}

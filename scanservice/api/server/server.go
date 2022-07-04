package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Atop-NMS-team/pgutils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	up "nms/api/v1/atopudpscan"
	ss "nms/api/v1/snmpscan"
	pb "nms/api/v1/scanservice"

	service "nms/scanservice/api/interface"

	"nms/scanservice/api/utils"
)

// server is used to implement
type server struct {
	pb.UnimplementedScanServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) StartScan(ctx context.Context, in *pb.StartScanRequest) (*pb.StartScanResponse, error) {
	log.Printf("parameter: %v", in)
	// -----------------------------------------------------------------------------------------------------
	// try to get watcher host for SNMP
	// -----------------------------------------------------------------------------------------------------
	sh := service.GetSnmpHost()
	if sh["status"] != 0 {
		return StartScanFailure(sh["message"].(string))
	}
	shh, ok := sh["response"].(string)
	if !ok {
		return StartScanFailure("could not change snmp host into string")
	}
	// -----------------------------------------------------------------------------------------------------
	// try to get watcher host for GWD
	// -----------------------------------------------------------------------------------------------------
	gh := service.GetGwdHost()
	if gh["status"] != 0 {
		return StartScanFailure(gh["message"].(string))
	}
	ghh, ok := gh["response"].(string)
	if !ok {
		return StartScanFailure("could not change gwd host into string")
	}
	// -----------------------------------------------------------------------------------------------------
	// try to connect database
	// -----------------------------------------------------------------------------------------------------
	cd := service.ConnectDB()
	if cd["status"] != 0 {
		return StartScanFailure(cd["message"].(string))
	}
	ccd, ok := cd["response"].(pgutils.IDBClient)
	if !ok {
		return StartScanFailure("could not change database connect into pgutils.IDBClient")
	}
	// -----------------------------------------------------------------------------------------------------
	// try to connect SNMP service
	// -----------------------------------------------------------------------------------------------------
	cs := service.ConnectSNMP(shh)
	if cs["status"] != 0 {
		return StartScanFailure(cs["message"].(string))
	}
	ccs := cs["response"].(ss.SnmpScanClient)
	if !ok {
		return StartScanFailure("could not change snmp client connect into SnmpScanClient")
	}
	// -----------------------------------------------------------------------------------------------------
	// try to connect GWD service
	// -----------------------------------------------------------------------------------------------------
	cg := service.ConnectGWD(ghh)
	if cg["status"] != 0 {
		return StartScanFailure(cg["message"].(string))
	}
	ccg := cg["response"].(up.GwdClient)
	if !ok {
		return StartScanFailure("could not change gwd client connect into GwdClient")
	}
	// -----------------------------------------------------------------------------------------------------
	// create session ID and insert into database
	// -----------------------------------------------------------------------------------------------------
	sessionId := time.Now().Format("20060102150405")
	create := service.CreateSession(ccd, sessionId)
	if create["status"] != 0 {
		return StartScanFailure(create["message"].(string))
	}
	// -----------------------------------------------------------------------------------------------------
	// start scan (SNMP)
	// -----------------------------------------------------------------------------------------------------
	ss := service.SnmpStartScan(ccs, in.Range, in.SnmpSettings, in.Oids, sessionId)
	// -----------------------------------------------------------------------------------------------------
	// start scan (GWD)
	// -----------------------------------------------------------------------------------------------------
	gs := service.GwdStartScan(ccg, in.ServerIp, sessionId)
	// -----------------------------------------------------------------------------------------------------
	// hanlder for scan fail
	// -----------------------------------------------------------------------------------------------------
	if ss["status"] != 0 || gs["status"] != 0 {
		query := service.QuerySession(ccd, sessionId)
		if query["status"] != 0 {
			return StartScanFailure(query["message"].(string))
		}
		session, ok := query["response"].(pgutils.DeviceSession)
		if !ok {
			return StartScanFailure("could not change session into pgutils.DeviceSession")
		}
		state := session.State
		if ss["status"] != 0 {
			state = utils.UpdateSessionState(state, "fail", "snmp")
		}
		if gs["status"] != 0 {
			state = utils.UpdateSessionState(state, "fail", "gwd")
		}
		update := service.UpdateState(ccd, session, state)
		if update["status"] != 0 {
			return StartScanFailure(update["message"].(string))
		}
	}
	// -----------------------------------------------------------------------------------------------------
	// success
	// -----------------------------------------------------------------------------------------------------
	return &pb.StartScanResponse{
		Info: &pb.SessionInfo{
			Success:   ss["status"] == 0 && gs["status"] == 0,
			SessionId: sessionId,
			Message:   ss["message"].(string) + "#" + gs["message"].(string),
		},
	}, nil
}

func (s *server) StopScan(ctx context.Context, in *pb.StopScanRequest) (*pb.StopScanResponse, error) {
	log.Printf("parameter: %v", in)
	// -----------------------------------------------------------------------------------------------------
	// success
	// -----------------------------------------------------------------------------------------------------
	return &pb.StopScanResponse{
		Info: &pb.Info{
			Success: true,
			Message: "",
		},
	}, nil
}

func (s *server) CheckStatus(ctx context.Context, in *pb.CheckStatusRequest) (*pb.CheckStatusResponse, error) {
	log.Printf("parameter: %v", in)
	// -----------------------------------------------------------------------------------------------------
	// try to connect database
	// -----------------------------------------------------------------------------------------------------
	cd := service.ConnectDB()
	if cd["status"] != 0 {
		return CheckStatusFailure(cd["message"].(string))
	}
	ccd, ok := cd["response"].(pgutils.IDBClient)
	if !ok {
		return CheckStatusFailure("could not change database connect into pgutils.IDBClient")
	}
	// -----------------------------------------------------------------------------------------------------
	// query session by session ID
	// -----------------------------------------------------------------------------------------------------
	query := service.QuerySession(ccd, in.SessionId)
	if query["status"] != 0 {
		return CheckStatusFailure(query["message"].(string))
	}
	session, ok := query["response"].(pgutils.DeviceSession)
	if !ok {
		return CheckStatusFailure("could not change session into pgutils.DeviceSession")
	}
	// -----------------------------------------------------------------------------------------------------
	// convert status (running | fail | success | notfound)
	// -----------------------------------------------------------------------------------------------------
	status := ""
	if strings.Contains(session.State, "fail") == true {
		status = "fail"
	} else if strings.Contains(session.State, "running") == true {
		status = "running"
	} else if strings.Contains(session.State, "success") == true {
		status = "success"
	} else {
		status = "notfound"
	}
	// -----------------------------------------------------------------------------------------------------
	// success
	// -----------------------------------------------------------------------------------------------------
	return &pb.CheckStatusResponse{
		Info: &pb.StatusInfo{
			Success: true,
			Status:  status,
			Message: "",
		},
	}, nil
}

func (s *server) GetResult(ctx context.Context, in *pb.GetResultRequest) (*pb.GetResultResponse, error) {
	log.Printf("parameter: %v", in)
	// -----------------------------------------------------------------------------------------------------
	// try to connect database
	// -----------------------------------------------------------------------------------------------------
	cd := service.ConnectDB()
	if cd["status"] != 0 {
		return GetResultFailure(cd["message"].(string))
	}
	ccd, ok := cd["response"].(pgutils.IDBClient)
	if !ok {
		return GetResultFailure("could not change database connect into pgutils.IDBClient")
	}
	// -----------------------------------------------------------------------------------------------------
	// query session by session ID
	// -----------------------------------------------------------------------------------------------------
	query1 := service.QuerySession(ccd, in.SessionId)
	if query1["status"] != 0 {
		return GetResultFailure(query1["message"].(string))
	}
	session, ok := query1["response"].(pgutils.DeviceSession)
	if !ok {
		return GetResultFailure("could not change session into pgutils.DeviceSession")
	}
	// -----------------------------------------------------------------------------------------------------
	// query result by session ID
	// -----------------------------------------------------------------------------------------------------
	query2 := service.QueryResult(ccd, in.SessionId)
	if query2["status"] != 0 {
		return GetResultFailure(query2["message"].(string))
	}
	result, ok := query2["response"].([]pgutils.DeviceResult)
	if !ok {
		return GetResultFailure("could not change result into []pgutils.DeviceResult")
	}
	// -----------------------------------------------------------------------------------------------------
	// format/pagination result
	// -----------------------------------------------------------------------------------------------------
	total := int32(len(result))
	content := Format(result, session)
	content = Pagination(content, in.Page, in.Size)
	// -----------------------------------------------------------------------------------------------------
	// success
	// -----------------------------------------------------------------------------------------------------
	return &pb.GetResultResponse{
		Info: &pb.Info{
			Success: true,
			Message: "",
		},
		Content: content,
		Page:    in.Page,
		Size:    in.Size,
		Total:   total,
	}, nil
}

func (s *server) GetLastResult(ctx context.Context, in *pb.GetLastResultRequest) (*pb.GetLastResultResponse, error) {
	log.Printf("parameter: %v", in)
	// -----------------------------------------------------------------------------------------------------
	// try to connect database
	// -----------------------------------------------------------------------------------------------------
	cd := service.ConnectDB()
	if cd["status"] != 0 {
		return GetLastResultFailure(cd["message"].(string))
	}
	ccd, ok := cd["response"].(pgutils.IDBClient)
	if !ok {
		return GetLastResultFailure("could not change database connect into pgutils.IDBClient")
	}
	// -----------------------------------------------------------------------------------------------------
	// query session by session ID
	// -----------------------------------------------------------------------------------------------------
	query1 := service.QueryLastSession(ccd)
	if query1["status"] != 0 {
		return GetLastResultFailure(query1["message"].(string))
	}
	session, ok := query1["response"].(pgutils.DeviceSession)
	if !ok {
		return GetLastResultFailure("could not change session into pgutils.DeviceSession")
	}
	// -----------------------------------------------------------------------------------------------------
	// query result by session ID
	// -----------------------------------------------------------------------------------------------------
	query2 := service.QueryResult(ccd, session.SessionID)
	if query2["status"] != 0 {
		return GetLastResultFailure(query2["message"].(string))
	}
	result, ok := query2["response"].([]pgutils.DeviceResult)
	if !ok {
		return GetLastResultFailure("could not change result into []pgutils.DeviceResult")
	}
	// -----------------------------------------------------------------------------------------------------
	// format/pagination result
	// -----------------------------------------------------------------------------------------------------
	total := int32(len(result))
	content := Format(result, session)
	content = Pagination(content, in.Page, in.Size)
	// -----------------------------------------------------------------------------------------------------
	// success
	// -----------------------------------------------------------------------------------------------------
	return &pb.GetLastResultResponse{
		Info: &pb.Info{
			Success: true,
			Message: "",
		},
		Content: content,
		Page:    in.Page,
		Size:    in.Size,
		Total:   total,
	}, nil
}

func Format(result []pgutils.DeviceResult, session pgutils.DeviceSession) []*pb.DeviceInfo {
	// create nil []*pb.DeviceInfo
	deviceInfo := []*pb.DeviceInfo{}
	// for loop
	for i := 0; i < len(result); i++ {
		// convert state
		status := "BOTH"
		states := strings.Split(session.State, "|")
		for i := 0; i < len(states); i++ {
			state := strings.Split(states[i], ":")
			if state[1] == "running" || state[1] == "fail" {
				status = "UNKNOW"
				break
			}
			if state[0] == "gwd" && state[1] != "success" {
				status = "SNMP"
			}
			if state[0] == "snmp" && state[1] != "success" {
				status = "GWD"
			}
		}
		// create DeviceInfo
		d := &pb.DeviceInfo{}
		d.SessionId = session.SessionID
		d.Model = result[i].Model
		d.MacAddress = result[i].MacAddress
		d.IpAddress = result[i].IpAddress
		d.Netmask = result[i].Netmask
		d.Gateway = result[i].Gateway
		d.Hostname = result[i].Hostname
		d.Kernel = result[i].Kernel
		d.Ap = result[i].Ap
		d.FirmwareVer = result[i].FirmwareVer
		d.Description = result[i].Description
		d.DeviceType = status
		d.ScanTime = session.LastUpdatedTime
		// append in to array
		deviceInfo = append(deviceInfo, d)
	}
	// return DeviceInfo array
	return deviceInfo
}

func StartScanFailure(message string) (*pb.StartScanResponse, error) {
	return &pb.StartScanResponse{
		Info: &pb.SessionInfo{
			Success:   false,
			SessionId: "",
			Message:   message,
		},
	}, nil
}

func CheckStatusFailure(message string) (*pb.CheckStatusResponse, error) {
	return &pb.CheckStatusResponse{
		Info: &pb.StatusInfo{
			Success: false,
			Status:  "",
			Message: message,
		},
	}, nil
}

func GetResultFailure(message string) (*pb.GetResultResponse, error) {
	return &pb.GetResultResponse{
		Info: &pb.Info{
			Success: false,
			Message: message,
		},
		Content: nil,
		Page:    0,
		Size:    0,
		Total:   0,
	}, nil
}

func GetLastResultFailure(message string) (*pb.GetLastResultResponse, error) {
	return &pb.GetLastResultResponse{
		Info: &pb.Info{
			Success: false,
			Message: message,
		},
		Content: nil,
		Page:    0,
		Size:    0,
		Total:   0,
	}, nil
}

func Pagination(result []*pb.DeviceInfo, page int32, size int32) []*pb.DeviceInfo {
	// create new array
	p := []*pb.DeviceInfo{}
	// start index
	start := (page - 1) * size
	// end index
	end := start + size
	// for loop
	for i := 0; i < len(result); i++ {
		if i >= int(start) && i < int(end) {
			p = append(p, result[i])
		}
	}
	// return pagination result
	return p
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func RunServer(grpcPort string, httpPort string) error {
	// normalize port
	gPort := utils.NormalizePort(grpcPort)

	// Create a listener on TCP port
	lis, err1 := net.Listen("tcp", gPort)
	if err1 != nil {
		log.Fatalln("Failed to listen:", err1)
		return err1
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pb.RegisterScanServiceServer(s, &server{})
	// Serve gRPC server
	log.Printf("Serving gRPC on 0.0.0.0%v \n", gPort)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err2 := grpc.DialContext(
		context.Background(),
		"0.0.0.0"+gPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err2 != nil {
		log.Fatalln("Failed to dial server:", err2)
		return err2
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err3 := pb.RegisterScanServiceHandler(context.Background(), gwmux, conn)
	if err3 != nil {
		log.Fatalln("Failed to register gateway:", err3)
		return err3
	}

	hPort := utils.NormalizePort(httpPort)

	gwServer := &http.Server{
		Addr:    hPort,
		Handler: cors(gwmux),
	}

	log.Printf("Serving gRPC-Gateway on http://0.0.0.0%v \n", hPort)
	log.Fatalln(gwServer.ListenAndServe())

	return nil
}

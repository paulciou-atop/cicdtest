package cmd

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"security/api/v1/security"
	"time"

	"github.com/bobbae/glog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewRootCrtCmd() *cobra.Command {

	runGetRootCrtCmd := &cobra.Command{
		Use:   "root",
		Short: "Get root.crt",
		Long:  "Get root.crt from Pki service",
		Run: func(cmd *cobra.Command, args []string) {
			serverip, err := cmd.Flags().GetString("server")
			if err != nil {
				glog.Error("Get serverip flag error, use localhost instead of ")
			}

			port, err := cmd.Flags().GetString("grpc-port")
			if err != nil {
				glog.Error("Get serverip flag error, use 8080 instead of ")
			}

			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			addr := net.JoinHostPort(serverip, port)
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("can not connect with server %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			client := security.NewPkiClient(conn)
			v := &security.EmptyParams{}
			stream, err := client.GetRootCrt(ctx, v)
			if err != nil {
				log.Fatalf("open stream error %v", err)
			}
			name := "root.crt"
			b := make(chan bool, 1)
			f := NewFile(name)
			go func() {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						if err := f.Save(); err != nil {
							b <- true
							return
						}
						b <- true
						return
					}
					if err != nil {
						log.Fatalf("cannot receive %v", err)
						b <- true
						return
					}

					if err := f.Write(resp.GetData()); err != nil {
						b <- true
						return
					}

				}
			}()

			<-b
		},
	}
	runGetRootCrtCmd.Flags().StringP("server", "s", "localhost", "server Ip")
	runGetRootCrtCmd.Flags().StringP("grpc-port", "g", "8080", "A grpc server listen port")

	return runGetRootCrtCmd
}
func NewSrvCrtCmd() *cobra.Command {

	runGetSrvCrtCmd := &cobra.Command{
		Use:   "srvcrt",
		Short: "Get srv.crt",
		Long:  "Get srv.crt from Pki sevice",
		Run: func(cmd *cobra.Command, args []string) {
			serverip, err := cmd.Flags().GetString("server")
			if err != nil {
				glog.Error("Get serverip flag error, use localhost instead of ")
			}

			port, err := cmd.Flags().GetString("grpc-port")
			if err != nil {
				glog.Error("Get serverip flag error, use 8080 instead of ")
			}
			addr := net.JoinHostPort(serverip, port)
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("can not connect with server %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			client := security.NewPkiClient(conn)
			v := &security.EmptyParams{}
			stream, err := client.GetSrvCrt(ctx, v)
			if err != nil {
				log.Fatalf("open stream error %v", err)
			}
			name := "srv.crt"
			b := make(chan bool, 1)
			f := NewFile(name)
			go func() {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						if err := f.Save(); err != nil {
							b <- true
							return
						}
						b <- true
						return
					}
					if err != nil {
						log.Fatalf("cannot receive %v", err)
						b <- true
						return
					}

					if err := f.Write(resp.GetData()); err != nil {
						b <- true
						return
					}

				}
			}()

			<-b
		},
	}
	runGetSrvCrtCmd.Flags().StringP("server", "s", "localhost", "server Ip")
	runGetSrvCrtCmd.Flags().StringP("grpc-port", "g", "8080", "A grpc server listen port")
	return runGetSrvCrtCmd
}

func NewSrvKeyCmd() *cobra.Command {

	runGetSrvKeyCmd := &cobra.Command{
		Use:   "srvkey",
		Short: "Get srv.key",
		Long:  "Get srv.key from Pki service",
		Run: func(cmd *cobra.Command, args []string) {
			serverip, err := cmd.Flags().GetString("server")
			if err != nil {
				glog.Error("Get serverip flag error, use localhost instead of ")
			}

			port, err := cmd.Flags().GetString("grpc-port")
			if err != nil {
				glog.Error("Get serverip flag error, use 8080 instead of ")
			}
			addr := net.JoinHostPort(serverip, port)
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("can not connect with server %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			client := security.NewPkiClient(conn)
			v := &security.EmptyParams{}
			stream, err := client.GetSrvKey(ctx, v)
			if err != nil {
				log.Fatalf("open stream error %v", err)
			}
			name := "srv.key"
			b := make(chan bool, 1)
			f := NewFile(name)
			go func() {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						if err := f.Save(); err != nil {
							b <- true
							return
						}
						b <- true
						return
					}
					if err != nil {
						log.Fatalf("cannot receive %v", err)
						b <- true
						return
					}

					if err := f.Write(resp.GetData()); err != nil {
						b <- true
						return
					}

				}
			}()

			<-b
		},
	}

	runGetSrvKeyCmd.Flags().StringP("server", "s", "localhost", "server Ip")
	runGetSrvKeyCmd.Flags().StringP("grpc-port", "g", "8080", "A grpc server listen port")
	return runGetSrvKeyCmd
}

type File struct {
	name   string
	buffer *bytes.Buffer
}

func NewFile(name string) *File {
	return &File{
		name:   name,
		buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}

func (f *File) Save() error {
	if err := ioutil.WriteFile(f.name, f.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

package client

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"security/api/v1/security"
	"security/config"

	"google.golang.org/grpc"
)

var port = config.GetgrpcPort()

func TestGetRootCrt(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, ":"+port, grpc.WithInsecure(), grpc.WithBlock())
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
}

func TestSrvCrt(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	conn, err := grpc.DialContext(ctx, ":"+port, grpc.WithInsecure(), grpc.WithBlock())
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
}

func TestSrvKey(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	conn, err := grpc.DialContext(ctx, ":"+port, grpc.WithInsecure(), grpc.WithBlock())
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

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

/*
  repo package, repository to effectively use resources in the project
	whithout making it global. Like db connection, service connection etc..
*/
package repo

import (
	"context"
	"nms/lib/pgutils"
	mq "nms/messaging"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type DBProvider interface {
	DB() pgutils.IClient
}

type MQProvider interface {
	MQ() mq.IMQClient
}

// GRPCConnetionProvider provide grpc connection for specific services,
type ObjectProvider interface {
	// Connection return nil if error
	Connection(service string) *grpc.ClientConn
	AddConnection(service string, conn *grpc.ClientConn)
}

type IRepo interface {
	DBProvider
	MQProvider
	ObjectProvider
	Close()
}

var repoinstance IRepo
var once sync.Once

var TIMEOUT = time.Second * 15

func GetRepo(ctx context.Context) (IRepo, error) {
	var err error
	once.Do(func() {
		repoinstance, err = newRepo(ctx)
	})

	return repoinstance, err
}

func newRepo(ctx context.Context) (IRepo, error) {
	r := &repo{db: nil, connections: map[string]*grpc.ClientConn{}}
	timeout, cancel := context.WithTimeout(ctx, TIMEOUT)
	r.db = dialDb(timeout)
	defer cancel()

	r.mq = mq.DialMessaging(timeout)

	return r, nil
}

type repo struct {
	db          pgutils.IClient
	mq          mq.IMQClient
	connections map[string]*grpc.ClientConn
}

func (r *repo) DB() pgutils.IClient {
	return r.db
}

func (r *repo) MQ() mq.IMQClient {
	return r.mq
}

func (r *repo) Close() {
	if r.db != nil {
		r.db.Close()
	}
	if r.mq != nil {
		r.mq.Close()
	}
	for _, s := range r.connections {
		s.Close()
	}

}

func (r *repo) Connection(service string) *grpc.ClientConn {
	c, ok := r.connections[service]
	if ok {
		return c
	}
	return nil
}

func (r *repo) AddConnection(service string, conn *grpc.ClientConn) {
	r.connections[service] = conn
}

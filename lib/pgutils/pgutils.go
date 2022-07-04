package pgutils

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

var (
	ErrNullDB       = errors.New("got null db")
	ErrDBNotSupport = errors.New("database not support")
)

type CreateTableOpt struct {
	IfNotExists bool
	Temp        bool
}

type ICreateTable interface {
	CreateTable(modle interface{}, opt ...CreateTableOpt) error
}

type QueryExpr struct {
	Expr  string
	Value interface{}
}

type Updater interface {
	// Update data's type should create schema first.
	// example: err := client.Update(MyStruct{Name:"austin"})
	Update(update interface{}) error
}

type Inserter interface {
	// Insert data's type should create schema first.
	// example: err := client.Insert(MyStruct{Name:"austin"})
	Insert(data interface{}) error
}

type Querier interface {
	// Query
	// results should be a slice, ex: []Author and db should have its shema
	/* example:
	var result = []devices{}
	err = c.Query(&result, "state = ?", "running")
	if err != nil {
		t.Error(err)
	}
	*/
	Query(results interface{}, queries ...QueryExpr) error
}
type IDataAccess interface {
	Querier
	Inserter
	Updater
}

// IClient
type IClient interface {
	ICreateTable
	IDataAccess
	GetDB() (*pg.DB, error)
	Close()
}

type NewClientOpt struct {
	Addr     string
	User     string
	Password string
	Database string
	Client   string
}

var defaultNewClientOpt = NewClientOpt{
	User:     "user",
	Password: "pass",
	Database: "nms",
	Client:   "postgres",
}

// NewClient new PostgreSQL client
func NewClient(opt ...NewClientOpt) (IClient, error) {
	var o NewClientOpt
	if len(opt) >= 1 {
		o = opt[0]
	} else {
		o = defaultNewClientOpt
	}
	switch o.Client {
	case "postgres":
		pgdb, err := connectPostgreDB(o)
		if err != nil {
			log.Errorf("connect postgres fail %v,use dummy client instead of.  ", err)
			return &DummyClient{}, err
		}

		return &PgClient{
			db: pgdb,
		}, nil
	default:
		return &DummyClient{}, fmt.Errorf("db: %s create client fail: %v, use dummy client instead of", o.Client, ErrDBNotSupport)
	}

}

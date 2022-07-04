package pgutils

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
)

type PgClient struct {
	db *pg.DB
}

func (c *PgClient) GetDB() (*pg.DB, error) {
	if c.db != nil {
		return c.db, nil
	}
	db, err := connectPostgreDB(defaultNewClientOpt)
	c.db = db
	return c.db, err

}

func (c *PgClient) Insert(data interface{}) error {
	_, err := c.db.Model(data).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (c *PgClient) Query(results interface{}, queries ...QueryExpr) error {
	query := c.db.Model(results)
	for _, q := range queries {
		query.Where(q.Expr, q.Value)
	}
	err := query.Select()
	if err != nil {
		return err
	}
	return nil
}

func (c *PgClient) Update(data interface{}) error {

	_, err := c.db.Model(data).WherePK().Update()
	if err != nil {
		return err
	}
	return nil
}

func (c *PgClient) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *PgClient) CreateTable(model interface{}, opt ...CreateTableOpt) error {
	if c.db == nil {
		log.Error("create table error ", ErrNullDB)
		return ErrNullDB
	}
	var createTableOpt = &orm.CreateTableOptions{}
	if len(opt) >= 1 {
		createTableOpt.IfNotExists = opt[0].IfNotExists
		createTableOpt.Temp = opt[0].Temp
	}
	err := c.db.Model(model).CreateTable(createTableOpt)

	return err
}

// connectPostgreDB
func connectPostgreDB(opt NewClientOpt) (*pg.DB, error) {
	url := getDBHost()

	pgdb := pg.Connect(&pg.Options{
		Addr:     url,
		User:     opt.User,
		Password: opt.Password,
		Database: opt.Database,
	})

	// check databas is up
	if err := pgdb.Ping(ctx); err != nil {
		log.Errorf("database not up err:%v", err)
		return nil, err
	}
	return pgdb, nil
}

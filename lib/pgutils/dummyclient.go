package pgutils

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type DummyClient struct {
}

func (c *DummyClient) GetDB() (*pg.DB, error) {

	return nil, fmt.Errorf("dummy client doesn't support GetDB()")

}

func (c *DummyClient) Insert(data interface{}) error {
	logrus.Infof("DUMMY- insert data:\n %v", data)
	return nil
}

func (c *DummyClient) Query(results interface{}, queries ...QueryExpr) error {
	logrus.Infof("DUMMY- querys ")
	for _, v := range queries {
		logrus.Infof(">> %v", v)
	}

	return nil
}

func (c *DummyClient) Update(data interface{}) error {
	logrus.Infof("DUMMY- update data:\n %v", data)

	return nil
}

func (c *DummyClient) Close() {
	logrus.Infof("DUMMY- close")
}

func (c *DummyClient) CreateTable(model interface{}, opt ...CreateTableOpt) error {
	logrus.Infof("DUMMY- create table model:\n")
	logrus.Infof("%v\n", model)
	return nil
}

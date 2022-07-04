/*
 This file consis of database related functions
*/
package repo

import (
	"context"
	"nms/lib/pgutils"
	"time"

	"github.com/sirupsen/logrus"
)

var retry = 3

func dialDb(ctx context.Context) pgutils.IClient {
	c := make(chan pgutils.IClient)
	logrus.Info("connect db ...")
	go func() {
		for {
			client, err := pgutils.NewClient()
			if err == nil {
				c <- client
				return
			}
			time.Sleep(time.Second * 3)
		}
	}()
	select {
	case <-ctx.Done():
		logrus.Info("cancel connect to db, use dummy client instead")
		return &pgutils.DummyClient{}
	case db := <-c:
		close(c)
		return db
	}
}

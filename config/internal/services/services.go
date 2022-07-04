package services

import (
	"nms/lib/repo"
	"nms/serviceswatcher/pkg/watcher"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ConfigerList = []string{"dummyconfiger", "atopudpscan", "snmpconfig"}

func InitServices(r repo.IRepo) error {
	// ctx, _ := GetGlobalContext()
	for _, c := range ConfigerList {
		service, err := watcher.GetServiceHostUrl(c)
		if err != nil {
			logrus.Info(err.Error())
		}
		if err == nil {
			//success
			conn, err := grpc.Dial(service, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}
			logrus.Infof("Found configer %s, running initial on %s...", c, service)
			r.AddConnection(c, conn)
		}
	}
	return nil
}

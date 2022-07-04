package device

import (
	"context"
	"nms/config/pkg/config"
	"sync"
)

type batchArg struct {
	dev     config.Device
	metrics []*config.ConfigMetric
}

type batchRet struct {
	sessionid  string
	deviceID   string
	devicePath string
	e          error
}

func batchconfig(ctx context.Context, c config.IConfig, configs []batchArg) <-chan batchRet {

	done := make(chan batchRet, len(configs))
	go func() {

		configDone := make(chan int)
		wg := new(sync.WaitGroup)
		wg.Add(len(configs))
		for _, v := range configs {
			go func(conf batchArg) {
				ret, err := c.Config(ctx, conf.dev, conf.metrics, true, configDone)
				if err != nil {
					done <- batchRet{e: err}
					goto Done
				}

				select {
				case <-configDone:
					done <- batchRet{sessionid: ret.Session.Id, deviceID: ret.Device.DeviceId, devicePath: ret.Device.DevicePath, e: nil}
					goto Done
				case <-ctx.Done():
					close(done)
				}
			Done:
				wg.Done()
			}(v)
		}
		wg.Wait()
		close(done)
	}()
	return done
}

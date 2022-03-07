/*
-- @Time : 2022/3/7 9:21
-- @Author : raoxiaoya
-- @Desc :
*/
package tracerhelper

import (
	"sync"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/phprao/go-skywalking.git/tracerhelper/util"
)

var tracer *go2sky.Tracer
var gcm util.GoroutineContextManager
var once sync.Once

func StartTracer(serviceAddr string, serviceName string) error {
	rp, err := reporter.NewGRPCReporter(serviceAddr, reporter.WithCheckInterval(time.Second))
	if err != nil {
		return err
	}
	tracer, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(rp))

	once.Do(func() {
		gcm = util.GoroutineContextManager{}
	})

	return nil
}

func GetTracer() *go2sky.Tracer {
	return tracer
}

func GetGcm() *util.GoroutineContextManager {
	return &gcm
}
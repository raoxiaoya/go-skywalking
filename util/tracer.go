/*
-- @Time : 2022/3/7 9:21
-- @Author : raoxiaoya
-- @Desc :
*/
package util

import (
	"sync"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

var tracer *go2sky.Tracer
var gcm GoroutineContextManager
var once sync.Once

func StartTracer() error {
	rp, err := reporter.NewGRPCReporter("192.168.2.44:11800", reporter.WithCheckInterval(time.Second))
	if err != nil {
		return err
	}
	tracer, err = go2sky.NewTracer("test-demo1", go2sky.WithReporter(rp))

	once.Do(func() {
		gcm = GoroutineContextManager{}
	})

	return nil
}

func GetTracer() *go2sky.Tracer {
	return tracer
}

func GetGcm() *GoroutineContextManager {
	return &gcm
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
	"github.com/phprao/go-skywalking.git/util"
)

var tr *go2sky.Tracer

func main() {
	r := gin.New()
	rp, err := reporter.NewGRPCReporter("192.168.2.44:11800", reporter.WithCheckInterval(time.Second))
	if err != nil {
		fmt.Println("create gosky reporter failed!")
		return
	}
	defer rp.Close()

	tr, err = go2sky.NewTracer("test-demo3", go2sky.WithReporter(rp))
	r.Use(util.Middleware(tr))
	r.GET("/test", test)
	r.Run(":7003")
}

func test(c *gin.Context) {
	result := make(map[string]interface{})
	result["code"] = 0
	result["msg"] = ""
	result["data"] = "test"
	c.JSON(http.StatusOK, result)
}
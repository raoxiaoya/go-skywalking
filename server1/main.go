package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phprao/go-skywalking.git/util"
)

func main() {
	if util.StartTracer() != nil {
		fmt.Println("create gosky reporter failed!")
	}
	r := gin.New()
	r.Use(util.Middleware())
	r.GET("/test", test)
	_ = r.Run(":7001")
}

func test(c *gin.Context) {
	util.Get("http://127.0.0.1:7002/test")

	result := make(map[string]interface{})
	result["code"] = 0
	result["msg"] = ""
	result["data"] = "test"
	c.JSON(http.StatusOK, result)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phprao/go-skywalking.git/model"
	"github.com/phprao/go-skywalking.git/tracerhelper"
	"github.com/phprao/go-skywalking.git/tracerhelper/ginagent"
)

func main() {
	if tracerhelper.StartTracer("192.168.2.44:11800", "test-demo3") != nil {
		fmt.Println("create gosky reporter failed!")
	}

	model.Setup()
	defer model.CloseAllDb()

	r := gin.New()
	r.Use(ginagent.Middleware())
	r.GET("/test", test)
	_ = r.Run(":7003")
}

func test(c *gin.Context) {
	model.Read5ScoreLogModel{}.GetId(1, 2)
	model.Read5WhiteListModel{}.GetId(1, 2)

	result := make(map[string]interface{})
	result["code"] = 0
	result["msg"] = ""
	result["data"] = "test"
	c.JSON(http.StatusOK, result)
}

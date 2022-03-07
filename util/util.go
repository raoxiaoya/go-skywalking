/*
-- @Time : 2022/3/2 10:54
-- @Author : raoxiaoya
-- @Desc :
*/
package util

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/gin-gonic/gin"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

func Get(link string) (response string, err error) {
	client := http.Client{Timeout: time.Second * 10}
	var reqest *http.Request
	reqest, err = http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}

	tracer = GetTracer()
	// perr 的作用就是Tag信息
	// operationName 就是调用名称，意思要明确
	url := reqest.URL
	operationName := url.Scheme + "://" + url.Host + url.Path
	ctx, _ := GetGcm().GetContext()
	span, err := tracer.CreateExitSpan(*ctx, operationName, url.Host, func(key, value string) error {
		reqest.Header.Set(key, value)
		return nil
	})
	if err != nil {
		return
	}
	span.SetComponent(componentIDGINHttpServer)
	span.Tag(go2sky.TagHTTPMethod, reqest.Method)
	span.Tag(go2sky.TagURL, link)
	span.SetSpanLayer(agentv3.SpanLayer_Http)

	resp, err := client.Do(reqest)
	if err != nil {
		return response, err
	}

	span.End()

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

const componentIDGINHttpServer = 5006

func Middleware() gin.HandlerFunc {
	tracer = GetTracer()
	if tracer == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		operationName := c.FullPath()
		span, ctx, err := tracer.CreateEntrySpan(c.Request.Context(), operationName, func(key string) (string, error) {
			return c.Request.Header.Get(key), nil
		})
		if err != nil {
			c.Next()
			return
		}
		span.SetComponent(componentIDGINHttpServer)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.SetSpanLayer(agentv3.SpanLayer_Http)

		// test log
		span.Log(time.Now(), "test log info")
		span.Error(time.Now(), "test log error")

		//c.Request = c.Request.WithContext(ctx)
		GetGcm().SetContext(&ctx)
		defer GetGcm().DelContext()

		c.Next()

		if len(c.Errors) > 0 {
			span.Error(time.Now(), c.Errors.String())
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
		span.End()
	}
}

/*
-- @Time : 2022/3/2 10:54
-- @Author : raoxiaoya
-- @Desc :
*/
package tracerhelper

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
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
	span.SetComponent(ComponentIDGINHttpServer)
	span.Tag(go2sky.TagHTTPMethod, reqest.Method)
	span.Tag(go2sky.TagURL, link)
	span.SetSpanLayer(agentv3.SpanLayer_Http)

	resp, err := client.Do(reqest)
	if err != nil {
		span.Error(time.Now(), err.Error())
	} else {
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(resp.StatusCode))
	}

	span.End()

	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
package request_generator

import (
	"github.com/x-armory/go-exception"
	"net/http"
	"net/url"
	"strings"
)

func ParseHttpRequest(request *Request) *http.Request {
	return GenRequest(request.Method, request.Url, request.Headers, request.Values)
}

// 根据http request所需参数组装http request，并设置默认header，避免被反爬
func GenRequest(method string, urlStr string, headers map[string][]string, values map[string][]string) *http.Request {
	request, err := http.NewRequest(method, urlStr, strings.NewReader(url.Values(values).Encode()))
	ex.AssertNoError(err)
	if headers != nil {
		request.Header = headers
	}
	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	}
	if request.Header.Get("Accept") == "" {
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	}
	if request.Header.Get("Accept-Language") == "" {
		request.Header.Set("Accept-Language", "zh-CN,zh")
	}
	if request.Header.Get("Accept-Encoding") == "" {
		request.Header.Set("Accept-Encoding", "")
	}
	if request.Header.Get("Pragma") == "" {
		request.Header.Set("Pragma", "no-cache")
	}
	if request.Header.Get("Cache-Control") == "" {
		request.Header.Set("Cache-Control", "no-cache")
	}
	if request.Header.Get("Upgrade-Insecure-Requests") == "" {
		request.Header.Set("Upgrade-Insecure-Requests", "1")
	}
	if request.Header.Get("Connection") == "" {
		request.Header.Set("Connection", "keep-alive")
	}
	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if request.Header.Get("Cache-control") == "" {
		request.Header.Set("Cache-control", "no-cache")
	}
	if request.Header.Get("pragma") == "" {
		request.Header.Set("pragma", "no-cache")
	}
	return request
}

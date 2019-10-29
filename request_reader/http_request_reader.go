package request_reader

import (
	"bytes"
	"encoding/json"
	"github.com/x-armory/go-crawler/util"
	"github.com/x-armory/go-exception"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var DefaultHttpRequestReader = &HttpRequestReader{
	Ignore404:  true,
	Client:     util.DefaultHttpClient,
	LogRequest: true,
	RetryTimes: 5,
}

type HttpRequestReader struct {
	Ignore404   bool
	Client      *http.Client
	LogRequest  bool
	LogResponse bool
	RetryTimes  int
}

func (r *HttpRequestReader) ReadRequest(req interface{}) io.Reader {
	request, ok := req.(*http.Request)
	if !ok {
		ex.Wrap("not supported parameter type").Throw()
	}
	var bodyBytes []byte
	if request.GetBody != nil {
		originBodyReader, e := request.GetBody()
		if e != nil {
			println("[WARN]", e.Error())
		} else {
			bodyBytes, _ = ioutil.ReadAll(originBodyReader)
		}
	}
	if r.LogRequest {
		// method, url
		println("[Method]", request.Method)
		println("[Url]", request.URL.String())
		// headers
		var bodyBytes []byte
		headerBytes, _ := json.Marshal(request.Header)
		println("[Headers]", string(headerBytes))
		// body
		println("[Body]", string(bodyBytes))
	}

	var response *http.Response
	var e error
	retryTimes := 0
	for true {
		newRequest, _ := http.NewRequest(request.Method, request.URL.String(), bytes.NewReader(bodyBytes))
		newRequest.Header = request.Header
		response, e = r.Client.Do(newRequest)
		if e == nil {
			break
		}
		retryTimes++
		if retryTimes > r.RetryTimes {
			println("[WARN]", e.Error(), retryTimes, "times, exit")
			break
		}
		println("[WARN]", e.Error(), retryTimes, "times, retry")
		time.Sleep(time.Second * time.Duration(retryTimes*retryTimes*retryTimes))
	}

	ex.AssertNoError(e, "http do request failed")
	ex.Assert(response.Body != nil, "response body is nil")
	defer response.Body.Close()
	ex.Assert(response.StatusCode < 500, "server error")
	if response.StatusCode == 404 {
		if !r.Ignore404 {
			panic(404)
		} else {
			return strings.NewReader("")
		}
	} else {
		ex.Assert(response.StatusCode < 400, response.StatusCode)
	}
	responseBytes, e := ioutil.ReadAll(response.Body)
	ex.AssertNoError(e, "read response error")
	if r.LogResponse {
		println("[Response] %s", string(responseBytes))
	}

	return bytes.NewReader(responseBytes)
}

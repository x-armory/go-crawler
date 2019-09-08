package request_reader

import (
	"github.com/x-armory/go-crawler/util"
	"github.com/x-armory/go-exception"
	"io"
	"net/http"
	"strings"
)

var DefaultHttpRequestReader = &HttpRequestReader{util.DefaultHttpClient}

type HttpRequestReader struct {
	Ignore404 bool
	Client    *http.Client
}

func (r *HttpRequestReader) ReadRequest(req interface{}) io.Reader {
	request, ok := req.(*http.Request)
	if !ok {
		ex.Wrap("not supported parameter type").Throw()
	}
	response, e := r.Client.Do(request)
	ex.AssertNoError(e, "http do request failed")
	ex.Assert(response.StatusCode < 500, "server error")
	if response.StatusCode == 404 {
		if !r.Ignore404 {
			panic("404")
		} else {
			return strings.NewReader("")
		}
	} else {
		ex.Assert(response.StatusCode < 400, response.Status)
	}
	return response.Body
}

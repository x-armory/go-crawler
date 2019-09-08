package request_reader

import (
	"github.com/x-armory/go-crawler/util"
	"github.com/x-armory/go-exception"
	"io"
	"net/http"
)

var DefaultHttpRequestReader = &HttpRequestReader{util.DefaultHttpClient}

type HttpRequestReader struct {
	Client *http.Client
}

func (r *HttpRequestReader) ReadRequest(req interface{}) io.Reader {
	request, ok := req.(*http.Request)
	if !ok {
		ex.Wrap("not supported parameter type").Throw()
	}
	response, e := r.Client.Do(request)
	ex.AssertNoError(e, "http do request failed")
	return response.Body
}

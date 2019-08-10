package crawler

import (
	"github.com/x-armory/go-crawler/util"
	"github.com/x-armory/go-exception"
	"github.com/x-armory/go-unmarshal-xpath"
	"gopkg.in/xmlpath.v2"
	"math/rand"
	"net/http"
	"time"
)

func NewXpathUnmarshaler(httpDelayMillisMin int, httpDelayMillisMax int, varStart int, varEnd int) *xpathUnmarshaler {
	ex.Assert(httpDelayMillisMin >= 0 && httpDelayMillisMax >= 0, "http delay should >=0")
	ex.Assert(httpDelayMillisMax >= httpDelayMillisMin, "httpDelayMillisMax should not less than httpDelayMillisMin")
	return &xpathUnmarshaler{
		httpDelayMillisMin,
		httpDelayMillisMax,
		xpath.XpathUnmarshaler{StartRow: varStart, EndRow: varEnd},
	}
}

type xpathUnmarshaler struct {
	httpDelayMillisMin int
	httpDelayMillisMax int
	xpath.XpathUnmarshaler
}

func (u *xpathUnmarshaler) Unmarshal(req *http.Request, target interface{}) {
	ex.Assert(req != nil, "request is nil")
	ex.Assert(target != nil, "business data is nil")
	response, e := util.DefaultHttpClient.Do(req)
	ex.AssertNoError(e)
	ex.Assert(response.Body != nil, "response body is nil")
	ex.Assert(response.StatusCode < 400, "http error code %d", response.StatusCode)
	defer response.Body.Close()
	node, e := xmlpath.ParseHTML(response.Body)
	ex.AssertNoError(e)
	ex.AssertNoError(u.XpathUnmarshal(node, target), "unmarshal failed")

	randInt := u.httpDelayMillisMax - u.httpDelayMillisMin
	if randInt > 0 {
		randInt = rand.Intn(randInt)
	}
	delayMillis := randInt + u.httpDelayMillisMin
	if delayMillis > 0 {
		time.Sleep(time.Millisecond * time.Duration(delayMillis))
	}
}

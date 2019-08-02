package crawler

import (
	"crypto/tls"
	"github.com/x-armory/go-exception"
	"github.com/x-armory/go-unmarshal-xpath"
	"gopkg.in/xmlpath.v2"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

var DefaultHttpClient *http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: time.Minute * 10,
	}
	proxy := os.Getenv("http_proxy")
	proxyUrl, err := url.Parse(proxy)
	if err == nil {
		transport.Proxy = http.ProxyURL(proxyUrl)
		println("use proxy", proxy)
	}
	DefaultHttpClient = &http.Client{Transport: transport, Timeout: time.Minute * 10}
}

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
	response, e := DefaultHttpClient.Do(req)
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

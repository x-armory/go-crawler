package crawler

import (
	"crypto/tls"
	"github.com/x-armory/go-exception"
	"github.com/x-armory/go-unmarshal-xpath"
	"gopkg.in/xmlpath.v2"
	"math/rand"
	"net/http"
	"time"
)

var DefaultHttpClient *http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: time.Minute * 10,
	}
	DefaultHttpClient = &http.Client{Transport: transport, Timeout: time.Minute * 10}
}

func NewXpathUnmarshaler(downloadIntervalMin int, downloadIntervalMax int, varStart int, varEnd int) *xpathUnmarshaler {
	ex.Assert(downloadIntervalMin >= 0 && downloadIntervalMax >= 0, "download interval should >=0")
	ex.Assert(downloadIntervalMax >= downloadIntervalMin, "downloadIntervalMax should not less than downloadIntervalMin")
	return &xpathUnmarshaler{
		downloadIntervalMin,
		downloadIntervalMax,
		xpath.XpathUnmarshaler{StartRow: varStart, EndRow: varEnd},
	}
}

type xpathUnmarshaler struct {
	downloadIntervalMin int
	downloadIntervalMax int
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
	u.XpathUnmarshal(node, target)

	randInt := u.downloadIntervalMax - u.downloadIntervalMin
	if randInt > 0 {
		randInt = rand.Intn(randInt)
	}
	sleepSeconds := randInt + u.downloadIntervalMin
	if sleepSeconds > 0 {
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}
}

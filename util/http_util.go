package util

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var DefaultHttpClient *http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSHandshakeTimeout: time.Minute * 10,
	}
	proxy := strings.TrimSpace(os.Getenv("http_proxy"))
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
			println("use proxy", proxy)
		}
	}
	DefaultHttpClient = &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 10,
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	println("[WARN] receive redirect response")
		//	if req.Response != nil && req.Response.Header != nil {
		//		location := req.Response.Header.Get("Location")
		//		println("[WARN] url", location)
		//	}
		//	return http.ErrUseLastResponse
		//},
	}
}

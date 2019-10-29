package util

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var DefaultHttpClient *http.Client
var DefaultHttpClientWithCookieJar *http.Client

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
	}
	DefaultHttpClientWithCookieJar = &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 10,
		Jar:       &cookiejar.Jar{},
	}
}

package crawler

import (
	"github.com/x-armory/go-exception"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SyncDuration int

const (
	Year SyncDuration = iota
	Month
	Day
)

func NewPeriodRequestGenerator(syncDuration SyncDuration, offset int, lastSyncTIme time.Time, parametersFunc PeriodRequestParametersFunc) *periodRequestGenerator {
	g := &periodRequestGenerator{syncDuration, offset, lastSyncTIme, parametersFunc}
	g.lastSyncTime = g.lastSyncTime.AddDate(0, 0, g.offset+1)
	switch g.duration {
	case Year:
		g.lastSyncTime = g.lastSyncTime.AddDate(-1, 1-int(g.lastSyncTime.Month()), 1-g.lastSyncTime.Day())
	case Month:
		g.lastSyncTime = g.lastSyncTime.AddDate(0, -1, 1-g.lastSyncTime.Day())
	case Day:
		g.lastSyncTime = g.lastSyncTime.AddDate(0, 0, -1)
	}
	return g
}

type periodRequestGenerator struct {
	duration       SyncDuration
	offset         int
	lastSyncTime   time.Time
	parametersFunc PeriodRequestParametersFunc
}

type PeriodRequestParametersFunc func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string)

func (g *periodRequestGenerator) GenRequest() *http.Request {
	return g.genRequest(g.parametersFunc(g.genNextPeriod()))
}

func (g *periodRequestGenerator) genNextPeriod() (start time.Time, end time.Time) {
	now := time.Now()
	if g.lastSyncTime.IsZero() || g.lastSyncTime.After(now) {
		return time.Time{}, time.Time{}
	}
	var end2 time.Time
	switch g.duration {
	case Year:
		g.lastSyncTime = g.lastSyncTime.AddDate(1, 0, 0)
		end2 = g.lastSyncTime.AddDate(1, 0, -1)
	case Month:
		g.lastSyncTime = g.lastSyncTime.AddDate(0, 1, 0)
		end2 = g.lastSyncTime.AddDate(0, 1, -1)
	case Day:
		g.lastSyncTime = g.lastSyncTime.AddDate(0, 0, 1)
		end2 = g.lastSyncTime
	}
	if g.lastSyncTime.After(now) {
		return time.Time{}, time.Time{}
	} else {
		return g.lastSyncTime, end2
	}
	return time.Time{}, time.Time{}
}

func (g *periodRequestGenerator) genRequest(method string, urlStr string, headers map[string][]string, values map[string][]string) *http.Request {
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
	return request
}

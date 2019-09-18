package request_generator

import (
	"github.com/x-armory/go-exception"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 适用于一个周期的所有数据放在一个url的场景；
// 需传入根据周期时间生成业务请求参数的方法；
// 已实现计算数据周期、转换http请求；
// SyncDuration：请求周期，支持年、月、日，默认日周期；
// Offset：天偏移量，有些数据在一个周期内分多次发布，有些数据会延迟发布，需要往前下载；
// LastTime：上次请求的时间点；
// ParametersFunc：业务定制的参数转换方法，根据一个周期的开始结束时间生成具体请求，日周期的开始结束时间一样；
// IgnoreWeekend：跳过周末，通常周末没有交易数据，只有日周期生效；
type DurationHttpRequestGenerator struct {
	Duration       SyncDuration
	Offset         int
	LastTime       time.Time
	IgnoreWeekend  bool
	ParametersFunc DurationRequestParametersFunc
	ready          bool
}
type SyncDuration int
type DurationRequestParametersFunc func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string)

const (
	Day SyncDuration = iota
	Month
	Year
)

func (g *DurationHttpRequestGenerator) GenRequest() interface{} {
	return g.genRequest(g.ParametersFunc(g.NextDuration()))
}

// 计算下一个周期的开始、结束时间
// 如果没有设置上次时间，或上次时间超过当前时间（异常情况），返回无效的时间段
// 如果计算出的下一个周期开始时间在当前时间之后，返回无效时间段
func (g *DurationHttpRequestGenerator) NextDuration() (start time.Time, end time.Time) {
	if !g.ready {
		// 加上偏移量后的下一天，后面再减掉一个周期，作为上次执行时间的开始时间
		g.LastTime = g.LastTime.AddDate(0, 0, g.Offset+1)
		// 格式化开始时间，减一个周期
		switch g.Duration {
		case Year: // 年周期的开始时间设为当年1月1日
			g.LastTime = g.LastTime.AddDate(-1, 1-int(g.LastTime.Month()), 1-g.LastTime.Day())
		case Month: // 月周期的开始时间为当月1日
			g.LastTime = g.LastTime.AddDate(0, -1, 1-g.LastTime.Day())
		case Day: //
			g.LastTime = g.LastTime.AddDate(0, 0, -1)
		}
		g.ready = true
	}
	ex.Assert(g.ParametersFunc != nil, "no ParametersFunc")

	now := time.Now()
	ex.Assert(!g.LastTime.IsZero(), "last sync time is zero")
	ex.Assert(!g.LastTime.After(now), "no more data")

	var end2 time.Time
	switch g.Duration {
	case Year:
		g.LastTime = g.LastTime.AddDate(1, 0, 0)
		end2 = g.LastTime.AddDate(1, 0, -1)
	case Month:
		g.LastTime = g.LastTime.AddDate(0, 1, 0)
		end2 = g.LastTime.AddDate(0, 1, -1)
	case Day:
		for true {
			g.LastTime = g.LastTime.AddDate(0, 0, 1)
			end2 = g.LastTime
			if g.LastTime.Weekday() >= 1 && g.LastTime.Weekday() <= 5 {
				break
			}
		}
	}
	ex.Assert(!g.LastTime.After(now), "no more data")

	return g.LastTime, end2
}

// 根据http request所需参数组装http request，并设置默认header，避免被反爬
func (g *DurationHttpRequestGenerator) genRequest(method string, urlStr string, headers map[string][]string, values map[string][]string) *http.Request {
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
	return request
}

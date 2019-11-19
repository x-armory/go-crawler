package request_generator

import (
	"time"
)

// 适用于一个周期的所有数据放在一个url的场景；
// 需传入根据周期时间生成业务请求参数的方法；
// 已实现计算数据周期、转换http请求；
// SyncDuration：请求周期，支持年、月、日，默认日周期；
// Offset：天偏移量，有些数据在一个周期内分多次发布，有些数据会延迟发布，需要往前下载；
// LastTime：上次请求的时间点；
// ParametersFunc：业务定制的参数转换方法，根据一个周期的开始结束时间生成具体请求，日周期的开始结束时间一样；
// DurationIsFinishedFunc：默认永远返回true；用于判断Duration是否完成，适用于每次请求有参数范围限制的场景，同一个时间段必须拆分成多次请求；
// IgnoreWeekend：跳过周末，通常周末没有交易数据，只有日周期生效；
type DurationHttpRequestGenerator struct {
	Duration               SyncDuration
	Offset                 int
	LastTime               time.Time
	IgnoreWeekend          bool
	ParametersFunc         DurationRequestParametersFunc
	DurationIsFinishedFunc func() bool
	ready                  bool
}
type SyncDuration int
type DurationRequestParametersFunc func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string)

const (
	Day SyncDuration = iota
	Month
	Year
)

func (g *DurationHttpRequestGenerator) GenRequest() interface{} {
	start, end := g.NextDuration()
	if start.IsZero() {
		println("[INFO] no more data")
		return nil
	}
	return GenRequest(g.ParametersFunc(start, end))
}

// 计算下一个周期的开始、结束时间
// 如果没有设置上次时间，或上次时间超过当前时间（异常情况），返回无效的时间段
// 如果计算出的下一个周期开始时间在当前时间之后，返回无效时间段
func (g *DurationHttpRequestGenerator) NextDuration() (start time.Time, end time.Time) {
	if !g.ready {
		if g.DurationIsFinishedFunc == nil {
			g.DurationIsFinishedFunc = func() bool {
				return true
			}
		}
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
	if g.ParametersFunc == nil {
		return time.Time{}, time.Time{}
	}
	now := time.Now()
	if g.LastTime.IsZero() || g.LastTime.After(now) {
		return time.Time{}, time.Time{}
	}
	var end2 time.Time
	switch g.Duration {
	case Year:
		if g.DurationIsFinishedFunc() {
			g.LastTime = g.LastTime.AddDate(1, 0, 0)
		}
		end2 = g.LastTime.AddDate(1, 0, -1)
	case Month:
		if g.DurationIsFinishedFunc() {
			g.LastTime = g.LastTime.AddDate(0, 1, 0)
		}
		end2 = g.LastTime.AddDate(0, 1, -1)
	case Day:
		if g.DurationIsFinishedFunc() {
			for true {
				g.LastTime = g.LastTime.AddDate(0, 0, 1)
				if !g.IgnoreWeekend || (g.LastTime.Weekday() >= 1 && g.LastTime.Weekday() <= 5) {
					break
				}
			}
		}
		end2 = g.LastTime
	}
	if g.LastTime.After(now) {
		return time.Time{}, time.Time{}
	} else {
		return g.LastTime, end2
	}
}

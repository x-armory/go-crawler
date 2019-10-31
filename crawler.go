package crawler

import (
	"bytes"
	"github.com/x-armory/go-exception"
	"github.com/x-armory/go-unmarshal/base"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

const UnmarshalError = "unmarshal failed"

// 同步执行RequestGenerator、RequestReader、DataUnmarshaler、DataProcessor方法
// 根据设定的间隔时间等待下次执行
// DataTarget：目标数据地址数组，元素可以是 *struct / []struct / []*struct；
// Ex：最后一个异常
// TimeInterval：控制RequestReader执行间隔，避免太快；
// TimeIntervalAddRand：随机延长TimeInterval，避免请求间隔太规律；
// DataUnmarshaler：执行反序列化，由Crawler的实现类选择具体的反序列化方法，以及处理过程，可以关闭写缓存，并选用ItemFilter来挨个处理元素；
// RequestGenerator：生成请求；
// RequestReader：读取请求返回内容；
// SkipUnmarshalError：true 无视UnmarshalError，打印异常后直接清空异常
// DurationFinally：可选项，每个间隔最终执行；
// Finally：可选项，最终执行；
type Crawler struct {
	DataTarget          []interface{}
	Ex                  *ex.ExceptionClass
	TimeInterval        time.Duration
	TimeIntervalAddRand time.Duration
	DataUnmarshaler     base.Unmarshaler
	RequestGenerator
	RequestReader
	SkipUnmarshalError bool
	DurationFinally    func(crawler *Crawler)
	Finally            func(crawler *Crawler)
}

// 生成请求参数
type RequestGenerator interface {
	GenRequest() interface{}
}

// 读取请求参数，生成Reader
type RequestReader interface {
	ReadRequest(req interface{}) (r io.Reader)
}

// 可选，反序列化完成后，处理全部数据
// 如果反序列化工具关闭了写数据，则无数据可处理
type DataProcessor interface {
	Process(target interface{})
}

func (c *Crawler) Start() {
	ex.Assert(c.RequestGenerator != nil, "RequestGenerator cannot be nil")
	ex.Assert(c.RequestReader != nil, "RequestReader cannot be nil")
	ex.Assert(c.DataUnmarshaler != nil, "DataUnmarshaler cannot be nil")
	rand.Seed(time.Now().UnixNano())

	for true {
		t1 := time.Now() //start time
		ex.Assert(len(c.DataTarget) > 0, "DataTarget cannot be empty")

		// read data
		var r io.Reader
		ex.Try(func() {
			req := c.GenRequest()
			ex.Assert(req != nil, "no more data")
			r = c.ReadRequest(req)
			ex.Assert(r != nil, "read data failed")
		}).Catch(func(err interface{}) {
			c.Ex = ex.Wrap(err)
		})

		// unmarshal data
		if c.Ex == nil {
			ex.Try(func() {
				// 如果目标对象超过1个，缓存io内容，用于以后读取
				var buf []byte
				if len(c.DataTarget) > 1 {
					buf, _ = ioutil.ReadAll(r)
				}
				for e := range c.DataTarget {
					var r2 = r
					if len(c.DataTarget) > 1 {
						r2 = bytes.NewReader(buf)
					}
					unmarshalErr := c.DataUnmarshaler.Unmarshal(r2, c.DataTarget[e])
					if unmarshalErr != nil {
						if c.SkipUnmarshalError {
							println("[WARN] skip unmarshal error:", unmarshalErr.Error())
						} else {
							ex.AssertNoError(unmarshalErr, UnmarshalError)
						}
					}
				}
			}).Catch(func(err interface{}) {
				c.Ex = ex.Wrap(err)
			})
		}

		// process data or err finally
		if c.DurationFinally != nil {
			ex.Try(func() {
				c.DurationFinally(c)
			}).Catch(func(err interface{}) {
				c.Ex = ex.Wrap(err)
			})
		}
		if c.Ex != nil {
			break
		}

		interval := time.Now().UnixNano() - t1.UnixNano() - int64(c.TimeInterval)
		if interval < 0 {
			interval = 0
		}
		if c.TimeIntervalAddRand > 0 {
			interval += rand.Int63n(int64(c.TimeIntervalAddRand))
		}
		if interval > 0 {
			time.Sleep(time.Duration(interval))
		}
	}

	if c.Finally != nil {
		ex.Try(func() {
			c.Finally(c)
		}).Catch(func(err interface{}) {
			ex.Wrap(err).PrintErrorStack()
		})
	} else {
		if c.Ex != nil {
			c.Ex.PrintErrorStack()
		}
	}
}

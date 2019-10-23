package crawler

import (
	"bytes"
	"github.com/x-armory/go-exception"
	"github.com/x-armory/go-unmarshal/base"
	"io"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

// 同步执行RequestGenerator、RequestReader、DataUnmarshaler、DataProcessor方法
// 根据设定的间隔时间等待下次执行
// DataTarget：目标数据地址数组，元素可以是 *struct / []struct / []*struct；
// Ex：最后一个异常
// TimeInterval：控制RequestReader执行间隔，避免太快；
// TimeIntervalAddRand：随机延长TimeInterval，避免请求间隔太规律；
// DataUnmarshaler：执行反序列化，由Crawler的实现类选择具体的反序列化方法，以及处理过程，可以关闭写缓存，并选用ItemFilter来挨个处理元素；
// RequestGenerator：生成请求；
// RequestReader：读取请求返回内容；
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
	DurationFinally func(crawler *Crawler)
	Finally         func(crawler *Crawler)
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
crawlerLoop:
	for true {
		ex.Assert(len(c.DataTarget) > 0, "DataTarget cannot be empty")
		// 用于控制请求间隔，最少间隔c.TimeInterval，如果请求处理时间超过c.TimeInterval，则处理时间为请求间隔
		// 然后随机再等待c.TimeIntervalAddRand
		// 确保请求不会很频繁、规律
		var bizFailedSig = make(chan bool)
		var bizFailed = false

		// 同步锁，用于等待业务执行结束
		wait := sync.WaitGroup{}
		wait.Add(1)

		// 异步执行业务代码
		go func() {
			defer func() {
				wait.Done()
				if bizFailed {
					bizFailedSig <- bizFailed
				}
			}()

			ex.Try(func() {
				req := c.GenRequest()
				if req == nil {
					c.Ex = ex.Wrap("no more data")
					bizFailed = true
					return
				}
				// 如果目标对象超过1个，缓存io内容，用于以后读取
				r := c.ReadRequest(req)
				var buf []byte
				if len(c.DataTarget) > 1 {
					buf, _ = ioutil.ReadAll(r)
				}
				for e := range c.DataTarget {
					var r2 = r
					if len(c.DataTarget) > 0 {
						r2 = bytes.NewReader(buf)
					}
					ex.AssertNoError(c.DataUnmarshaler.Unmarshal(r2, c.DataTarget[e]), "unmarshal failed")
				}
			}).Catch(func(err interface{}) {
				c.Ex = ex.Wrap(err)
			})

			ex.Try(func() {
				if c.DurationFinally != nil { // 如果定义了DurationFinally，将执行异常交给DurationFinally处理
					c.DurationFinally(c)
				} else {
					if c.Ex != nil { // 如果没定义DurationFinally，且出现了执行异常，直接抛出
						c.Ex.Throw()
					}
				}
			}).Catch(func(err interface{}) {
				c.Ex = ex.Wrap(err)
				bizFailed = true
			})

			// 如果执行出现了异常，且需要忽略，需要定义DurationFinally，并吃掉执行异常
		}()

		// 等待执行异常信号，或者间隔超时
		// 如果执行成功，则收不到异常信号，只能收到超时信号，然后开始下一轮
		// 如果执行失败，且执行时间大于间隔时间，会提前收到超时信号，进入后续等待执行并判断流程
		// 如果执行失败，且执行时间小于间隔时间，会及时收到异常信号并退出
		select {
		case <-bizFailedSig: // 收到退出信号
			break crawlerLoop
		case <-time.After(c.TimeInterval): // 执行到达间隔时间，此后可能执行成功，也可能失败
		}

		// 等待业务执行结束，此时一定是到达了间隔超时时间
		// 如果执行失败，会发出退出信号
		// 如果执行成功，不会发出退出信号，如果尝试读取会无限等待
		// 因此下一步使用异常状态来判断
		wait.Wait()
		if bizFailed {
			break crawlerLoop
		}
		close(bizFailedSig)

		// 如果设置了随机额外间隔时间，继续等待
		if c.TimeIntervalAddRand > 0 {
			time.Sleep(time.Duration(rand.Int63n(int64(c.TimeIntervalAddRand))))
		}
	}

	ex.Try(func() {
		if c.Finally != nil {
			c.Finally(c)
		}
	}).Catch(func(err interface{}) {
		ex.Wrap(err).PrintErrorStack()
	})
}

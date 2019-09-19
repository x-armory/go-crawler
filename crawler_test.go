package crawler

import (
	"fmt"
	"github.com/x-armory/go-exception"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestCrawler_Start(t *testing.T) {
	var data []Dto
	crawler := Crawler{
		DataTargetGetter: func() interface{} {
			return data
		},
		TimeInterval:        time.Second * 2,
		TimeIntervalAddRand: time.Duration(0),
		RequestGenerator:    &TestCrawler_RequestGenerator{},
		RequestReader:       &TestCrawler_RequestReader{},
		DataUnmarshaler:     &TestCrawler_DataUnmarshaler{},
		DurationFinally: func(data interface{}, err *ex.ExceptionClass) {
			if err != nil {
				err.PrintErrorStack()
			} else {
				println(fmt.Sprintf("%v", data))
			}
		},
		Finally: func(data interface{}, err *ex.ExceptionClass) {
			if err != nil {
				err.PrintErrorStack()
			} else {
				println(fmt.Sprintf("%v", data))
			}
		},
	}
	crawler.Start()
}

type Dto struct {
}
type TestCrawler_RequestGenerator struct {
}
type TestCrawler_RequestReader struct {
}
type TestCrawler_DataUnmarshaler struct {
}

var i = 1

func (g *TestCrawler_RequestGenerator) GenRequest() interface{} {
	req := fmt.Sprintf("req:%v", rand.Int63())
	i++
	if i == 5 {
		println("quit when i =", i)
		return nil
		//panic(req)
	}
	println(fmt.Sprintf("%v GenRequest(), return %s", time.Now(), req))
	return req
}
func (r *TestCrawler_RequestReader) ReadRequest(req interface{}) io.Reader {
	println(fmt.Sprintf("%v ReadRequest(%v)", time.Now(), req))
	return strings.NewReader(req.(string))
}
func (m *TestCrawler_DataUnmarshaler) Unmarshal(r io.Reader, target interface{}) error {
	println(fmt.Sprintf("%v Unmarshal(%T, %T)", time.Now(), r, target))
	return nil
}

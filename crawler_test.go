package crawler

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestCrawler_Start(t *testing.T) {
	var data []Dto
	crawler := Crawler{
		DataTarget:          data,
		TimeInterval:        time.Second * 2,
		TimeIntervalAddRand: time.Duration(0),
		RequestGenerator:    &TestCrawler_RequestGenerator{},
		RequestReader:       &TestCrawler_RequestReader{},
		DataUnmarshaler:     &TestCrawler_DataUnmarshaler{},
		DataProcessor:       &TestCrawler_DataProcessor{},
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
type TestCrawler_DataProcessor struct {
}

var i = 1

func (g *TestCrawler_RequestGenerator) GenRequest() interface{} {
	req := fmt.Sprintf("req:%v", rand.Int63())
	i++
	if i == 5 {
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
func (m *TestCrawler_DataUnmarshaler) Unmarshal(r io.Reader, target interface{}) {
	println(fmt.Sprintf("%v Unmarshal(%T, %T)", time.Now(), r, target))
}
func (p *TestCrawler_DataProcessor) Process(target interface{}) {
	println(fmt.Sprintf("%v Process(%T)", time.Now(), target))
	println("==============")
}

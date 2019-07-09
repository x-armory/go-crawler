package main

import (
	"bytes"
	"fmt"
	"github.com/x-armory/go-crawler"
	"github.com/x-armory/go-exception"
	"net/url"
	"time"
)

func main() {
	NewBaiduCrawler().Start()
}

// custom crawler business definition
func NewBaiduCrawler() *crawler.Crawler {
	lastSyncTime := time.Now().AddDate(0, 0, -3)
	return &crawler.Crawler{
		Business:         &BaiduBusiness{},
		RequestGenerator: crawler.NewPeriodRequestGenerator(crawler.Day, 0, lastSyncTime, getRequestParametersFunc()),
		DataUnmarshaler:  crawler.NewXpathUnmarshaler(0, 0, 1, -1),
		Notification:     crawler.NewCombinedNotification(&dingNotification{}),
	}
}

type BaiduData struct {
	Title string `xpath:"//*[@id='%d']/h3/a"`
	Desc  string `xpath:"//*[@id='%d']/div[1]/text()"`
}

type BaiduBusiness struct {
	data   []BaiduData
	report bytes.Buffer
	count  int
}

func (b *BaiduBusiness) NewPeriodData() interface{} {
	b.data = []BaiduData{}
	return &b.data
}

func (b *BaiduBusiness) ProcessPeriodData() {
	for _, d := range b.data {
		b.count++
		fmt.Printf("%+v\n", d)
	}
}

func (b *BaiduBusiness) GenReport() string {
	println("total count", b.count)
	return b.report.String()
}

func getRequestParametersFunc() crawler.PeriodRequestParametersFunc {
	return func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string) {
		ex.Assert(!start.IsZero(), ex.Exception(crawler.NoMoreDataException, "", nil))
		date := start.Format("2006-01-02")
		println("sync date", date)
		encode := url.Values(map[string][]string{"wd": {date}}).Encode()
		return "GET", "https://www.baidu.com/s?" + encode, nil, nil
	}
}

type dingNotification struct {
}

func (n *dingNotification) Send(s string) {
	println("ding:", s)
}

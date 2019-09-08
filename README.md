# go-crawler
get data easily to custom struct using annotations, 


##### sample:

```go
package main

import (
	zip2 "archive/zip"
	"github.com/x-armory/go-crawler"
	"github.com/x-armory/go-crawler/request_generator"
	"github.com/x-armory/go-crawler/request_reader"
	"github.com/x-armory/go-unmarshal/base"
	"github.com/x-armory/go-unmarshal/zip"
	"time"
)

type ShangHaiQiHuo struct {
	Name   string    `xm:"excel://sheet[0]/row[3:]/col[0] pattern='\\w+'"`
	Date   time.Time `xm:"excel://sheet[0]/row[3:]/col[1] pattern='\\d{8}' format='20060102' timezone='Asia/Shanghai'"`
	Open   int       `xm:"excel://sheet[0]/row[3:]/col[4] pattern='\\d+'"`
	High   int       `xm:"excel://sheet[0]/row[3:]/col[5] pattern='\\d+'"`
	Low    int       `xm:"excel://sheet[0]/row[3:]/col[6] pattern='\\d+'"`
	Close  int       `xm:"excel://sheet[0]/row[3:]/col[7] pattern='\\d+'"`
	Close2 int       `xm:"excel://sheet[0]/row[3:]/col[8] pattern='\\d+'"`
	Vol    int       `xm:"excel://sheet[0]/row[3:]/col[11] pattern='\\d+'"`
	Amount int       `xm:"excel://sheet[0]/row[3:]/col[12] pattern='\\d+'"`
	Cang   int       `xm:"excel://sheet[0]/row[3:]/col[13] pattern='\\d+'"`
}

func main() {
	c := &crawler.Crawler{
		DataTarget:          []*ShangHaiQiHuo{},
		TimeInterval:        time.Second,
		TimeIntervalAddRand: time.Second * 2,
		RequestReader:       request_reader.DefaultHttpRequestReader,
		RequestGenerator: &request_generator.DurationHttpRequestGenerator{
			Duration:      request_generator.Year,
			IgnoreWeekend: true,
			LastTime:      time.Now(),
			ParametersFunc: func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string) {

			},
		},
		DataUnmarshaler: &zip.Unmarshaler{
			Charset: "gbk",
			FileFilters: []zip.FileFilter{
				func(fileIndex int, file *zip2.File) bool {

				},
			},
			DataLoader: base.DataLoader{
				ItemFilters: []base.ItemFilter{
					func(item interface{}, vars *base.Vars) (flow base.FlowControl, deep int) {

					},
				},
			},
		},
		Finally: func() {

		},
	}
	c.Start()
}
```
package crawler

import (
	"github.com/x-armory/go-exception"
	"net/http"
)

type Crawler struct {
	RequestGenerator
	DataUnmarshaler
	Notification
	Business
}

func (c Crawler) Start() {
	ex.Assert(c.RequestGenerator != nil, "request generator is nil")
	ex.Assert(c.DataUnmarshaler != nil, "data unmarshaler is nil")
	ex.Assert(c.Notification != nil, "notification is nil")
	ex.Assert(c.Business != nil, "business definition is nil")
	ex.Try(func() {
		for true {
			c.Unmarshal(c.GenRequest(), c.NewPeriodData())
			c.ProcessPeriodData()
		}
	}).SafeCatch(func(err interface{}) {
		if e := ex.Wrap(err); e.Code() != NoMoreDataException {
			e.Throw()
		}
	})

	c.Send(c.GenReport())
}

type RequestGenerator interface {
	GenRequest() (req *http.Request)
}
type DataUnmarshaler interface {
	Unmarshal(req *http.Request, target interface{})
}
type Business interface {
	NewPeriodData() interface{}
	ProcessPeriodData()
	GenReport() string
}
type Notification interface {
	Send(n string)
}

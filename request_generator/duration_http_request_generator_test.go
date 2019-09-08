package request_generator

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestDurationHttpRequestGenerator_Gen_Daily_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-02")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-02")
	generator := &DurationHttpRequestGenerator{Duration: Day, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Daily_Weekend(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-04")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-07")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-07")
	generator := &DurationHttpRequestGenerator{Duration: Day, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Daily_Offset_Pre_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-01")
	generator := &DurationHttpRequestGenerator{Duration: Day, Offset: -1, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}
func TestDurationHttpRequestGenerator_Gen_Daily_Offset_Post_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-03")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-03")
	generator := &DurationHttpRequestGenerator{Duration: Day, Offset: 1, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Daily_LastDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-31")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-02-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-02-01")
	generator := &DurationHttpRequestGenerator{Duration: Day, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Daily_Future(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2819-01-31")
	generator := &DurationHttpRequestGenerator{Duration: Day, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start.IsZero(), true)
	assert.Equal(t, end.IsZero(), true)
}

func TestDurationHttpRequestGenerator_Gen_Monthly_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-07-31")
	generator := &DurationHttpRequestGenerator{Duration: Month, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Monthly_Offset_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-06-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-06-30")
	generator := &DurationHttpRequestGenerator{Duration: Month, Offset: -2, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Monthly_LastDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-06-30")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-07-31")
	generator := &DurationHttpRequestGenerator{Duration: Month, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Monthly_Future(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2819-07-01")
	generator := &DurationHttpRequestGenerator{Duration: Month, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start.IsZero(), true)
	assert.Equal(t, end.IsZero(), true)
}

func TestDurationHttpRequestGenerator_Gen_Yearly_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-12-31")
	generator := &DurationHttpRequestGenerator{Duration: Year, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Yearly_Offset_FirstDayOfYear(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2018-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2018-12-31")
	generator := &DurationHttpRequestGenerator{Duration: Year, Offset: -2, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Yearly_LastDayOfYear(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2018-12-31")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-12-31")
	generator := &DurationHttpRequestGenerator{Duration: Year, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestDurationHttpRequestGenerator_Gen_Yearly_Future(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2819-07-01")
	generator := &DurationHttpRequestGenerator{Duration: Year, Offset: 0, LastTime: lastSyncDate, IgnoreWeekend: true, ParametersFunc: getDurationHttpRequestParametersFunc()}
	start, end := generator.NextDuration()
	assert.Equal(t, start.IsZero(), true)
	assert.Equal(t, end.IsZero(), true)
}

func getDurationHttpRequestParametersFunc() DurationRequestParametersFunc {
	return func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string) {
		return "", "", nil, nil
	}
}

package crawler

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestPeriodRequestGenerator_Gen_Daily_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-02")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-02")
	generator := NewPeriodRequestGenerator(Day, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Daily_Weekend(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-04")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-07")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-07")
	generator := NewPeriodRequestGenerator(Day, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Daily_Offset_Pre_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-01")
	generator := NewPeriodRequestGenerator(Day, -1, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}
func TestPeriodRequestGenerator_Gen_Daily_Offset_Post_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-03")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-01-03")
	generator := NewPeriodRequestGenerator(Day, 1, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Daily_LastDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-31")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-02-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-02-01")
	generator := NewPeriodRequestGenerator(Day, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Daily_Future(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2819-01-31")
	generator := NewPeriodRequestGenerator(Day, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start.IsZero(), true)
	assert.Equal(t, end.IsZero(), true)
}

func TestPeriodRequestGenerator_Gen_Monthly_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-07-31")
	generator := NewPeriodRequestGenerator(Month, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Monthly_Offset_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-06-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-06-30")
	generator := NewPeriodRequestGenerator(Month, -2, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Monthly_LastDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-06-30")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-07-31")
	generator := NewPeriodRequestGenerator(Month, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Monthly_Future(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2819-07-01")
	generator := NewPeriodRequestGenerator(Month, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start.IsZero(), true)
	assert.Equal(t, end.IsZero(), true)
}

func TestPeriodRequestGenerator_Gen_Yearly_FirstDayOfMonth(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-07-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-12-31")
	generator := NewPeriodRequestGenerator(Year, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Yearly_Offset_FirstDayOfYear(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedStart, _ := time.Parse("2006-01-02", "2018-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2018-12-31")
	generator := NewPeriodRequestGenerator(Year, -2, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Yearly_LastDayOfYear(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2018-12-31")
	exceptedStart, _ := time.Parse("2006-01-02", "2019-01-01")
	exceptedEnd, _ := time.Parse("2006-01-02", "2019-12-31")
	generator := NewPeriodRequestGenerator(Year, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start, exceptedStart)
	assert.Equal(t, end, exceptedEnd)
}

func TestPeriodRequestGenerator_Gen_Yearly_Future(t *testing.T) {
	lastSyncDate, _ := time.Parse("2006-01-02", "2819-07-01")
	generator := NewPeriodRequestGenerator(Year, 0, lastSyncDate, getPeriodRequestParametersFunc(), true)
	start, end := generator.genNextPeriod()
	assert.Equal(t, start.IsZero(), true)
	assert.Equal(t, end.IsZero(), true)
}

func getPeriodRequestParametersFunc() PeriodRequestParametersFunc {
	return func(start time.Time, end time.Time) (method string, urlStr string, headers map[string][]string, values map[string][]string) {
		return "", "", nil, nil
	}
}

package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TT struct {
	x  int `expect:"8"`
	y  int `expect:"0"`
	z  int
	xf float64 `expect:"8"`
	yf float32 `expect:"0"`
	zf float32
	xb bool   `expect:"true"`
	yb bool   `expect:"false"`
	xs string `expect:"abc"`
	xy string `expect:""`
}

func TestGetEqualTagFields(t *testing.T) {
	tt := TT{x: 8, xf: 8.000, xs: "abc"}
	fields := GetEqualTagFields(&tt, "expect")
	fmt.Printf("%+v\n", fields)
	assert.Equal(t, fields, []string{"x", "y", "xf", "yf", "yb", "xs", "xy"})
}

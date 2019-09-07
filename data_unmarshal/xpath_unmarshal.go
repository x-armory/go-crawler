package data_unmarshal

import (
	"github.com/x-armory/go-exception"
	"github.com/x-armory/go-unmarshal-xpath"
	"gopkg.in/xmlpath.v2"
	"io"
)

func NewXpathUnmarshaler(varStart int, varEnd int) *xpathUnmarshaler {
	return &xpathUnmarshaler{
		xpath.XpathUnmarshaler{StartRow: varStart, EndRow: varEnd},
	}
}

type xpathUnmarshaler struct {
	xpath.XpathUnmarshaler
}

func (u *xpathUnmarshaler) Unmarshal(r io.Reader, target interface{}) {
	node, e := xmlpath.ParseHTML(r)
	ex.AssertNoError(e)
	ex.AssertNoError(u.XpathUnmarshal(node, target), "unmarshal failed")
}

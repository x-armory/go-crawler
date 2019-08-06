package util

import (
	"github.com/x-armory/go-exception"
	"reflect"
	"strconv"
)

// 检查所有字段，如果字段值和tag指定的值一样，则返回字段名列表
func GetEqualTagFields(v interface{}, tag string) []string {
	o := reflect.Indirect(reflect.ValueOf(v))
	T := reflect.TypeOf(v)
	for T.Kind() == reflect.Ptr {
		T = T.Elem()
	}
	var fields []string
	for i := 0; i < T.NumField(); i++ {
		field := T.Field(i)
		xormOmit, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		vField := o.FieldByName(field.Name)
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			omitV, ok := strconv.ParseInt(xormOmit, 10, 64)
			ex.AssertNoError(ok)
			if vField.Int() == omitV {
				fields = append(fields, field.Name)
			}
		case reflect.Float32, reflect.Float64:
			omitV, ok := strconv.ParseFloat(xormOmit, 64)
			ex.AssertNoError(ok)
			if vField.Float() == omitV {
				fields = append(fields, field.Name)
			}
		}
	}
	return fields
}

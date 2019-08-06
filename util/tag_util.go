package util

import (
	"github.com/x-armory/go-exception"
	"reflect"
	"strconv"
)

// 检查所有字段，如果字段值和tag指定的值一样，则返回字段名列表
// 目前支持类型：所有int、所有float、string、bool
func GetEqualTagFields(v interface{}, tag string) []string {
	o := reflect.Indirect(reflect.ValueOf(v))
	T := reflect.TypeOf(v)
	for T.Kind() == reflect.Ptr {
		T = T.Elem()
	}
	var fields []string
	for i := 0; i < T.NumField(); i++ {
		field := T.Field(i)
		tagValue, err := field.Tag.Lookup(tag)
		if !err {
			continue
		}
		vField := o.FieldByName(field.Name)
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tagRealValue, err := strconv.ParseInt(tagValue, 10, 64)
			ex.AssertNoError(err)
			if vField.Int() == tagRealValue {
				fields = append(fields, field.Name)
			}
		case reflect.Float32, reflect.Float64:
			tagRealValue, err := strconv.ParseFloat(tagValue, 64)
			ex.AssertNoError(err)
			if vField.Float() == tagRealValue {
				fields = append(fields, field.Name)
			}
		case reflect.String:
			if vField.String() == tagValue {
				fields = append(fields, field.Name)
			}
		case reflect.Bool:
			tagRealValue, err := strconv.ParseBool(tagValue)
			ex.AssertNoError(err)
			if vField.Bool() == tagRealValue {
				fields = append(fields, field.Name)
			}
		}
	}
	return fields
}

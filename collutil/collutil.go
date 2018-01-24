package collutil

import (
	"fmt"
	"reflect"

	"github.com/ghettovoice/gosip/log"
)

func Keys(input interface{}) []interface{} {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Map {
		log.Panicf("input type not a map: %v", v)
	}
	keys := make([]interface{}, 0, len(v.MapKeys()))
	for _, k := range v.MapKeys() {
		keys = append(keys, k.Interface())
	}
	return keys
}

func Vals(input interface{}) []interface{} {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Map {
		log.Panicf("input type not a map: %v", v)
	}
	vals := make([]interface{}, 0, len(v.MapKeys()))
	for _, k := range v.MapKeys() {
		vals = append(vals, v.MapIndex(k).Interface())
	}
	return vals
}

func ToSliceOfIface(input interface{}) []interface{} {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		log.Panicf("input type not a slice or map: %v", v)
	}
	slc := make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		slc = append(slc, v.Index(i).Interface())
	}
	return slc
}

func ToSliceOfString(input interface{}) []string {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		log.Panicf("input type not a slice or map: %v", v)
	}
	slc := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		slc = append(slc, fmt.Sprintf("%v", v.Index(i).Interface()))
	}
	return slc
}

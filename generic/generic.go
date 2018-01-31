// Package generic provides a set of generic functions to work with collections (arrays, slices, maps and chans).
package generic

import (
	"fmt"
	"reflect"
)

// Iteratee represents function that is called for each element of an collection.
// Function should have one, two or three input arguments and one or two output argument.
// The second iteratee output always expected to be error, that indicates early break of the loop.
//
// There are several common iteratee patterns expected:
//
// 1. Map-like iteratee with one or two inputs of any type and one or two outputs. The second output expected to be of
// error type.
//		// pseudo-code
//		func(val Any [, key Any]) (out Any [, err error])
//		// examples
//		func(val int, key string) int { return val }
//		func(val string) (int, error) { return strconv.Atoi(val) }
//
// 2. Filter-like iteratee with one or two inputs of any type and one or two outputs. The first output expected to be
// of bool type. The second output expected to be of error type.
//		// pseudo-code
//		func(val Any [, key Any]) (ok bool [, err error])
//		// examples
//		func(val int, key string) int { return val }
//		func(val string) (int, error) { return strconv.Atoi(val) }
//
// 3. Reduce-like iteratee with two or three inputs and one or two outputs. The second output expected to be
// of error type.
//		// pseudo-code
//		func(res Any, val Any [, key Any]) (out Any [, err error])
//		// examples
//		func(all []string, val string, key int) (string, error) { return all + fmt.Sprintf(" -> %d = %s", key, val)), nil }
//		func(all []float32, val float32) []float32 { return append(all, val * 2) }
type Iteratee interface{}

// Common iteratees written in pseudo-code.
const (
	PseudoMapLikeIteratee    = "func(val Any [, key Any]) (out Any [, err error])"
	PseudoFilterLikeIteratee = "func(val Any [, key Any]) (ok bool [, err error])"
	PseudoReduceLikeIteratee = "func(res Any, val Any [, key Any]) (out Any [, err error])"
)

// Collection represents a slice, array or map.
type Collection interface{}

// Pointer represents a pointer to result.
type Pointer interface{}

// Any represents value of any type.
type Any interface{}

var collectionKinds = []reflect.Kind{reflect.Array, reflect.Slice, reflect.Map, reflect.Chan}
var collectionPtrKinds []string

func init() {
	for _, kind := range collectionKinds {
		collectionPtrKinds = append(collectionPtrKinds, fmt.Sprintf("*%s", kind))
	}
}

// CollectionKinds returns slice of supported collection kinds.
func CollectionKinds() []reflect.Kind {
	return append([]reflect.Kind{}, collectionKinds...)
}

// CollectionPtrKinds returns list of supported collection pointers.
func CollectionPtrKinds() []string {
	return append([]string{}, collectionPtrKinds...)
}

//func Each(input interface{}, fn Iteratee) {
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//		for _, key := range v.MapKeys() {
//			if fn(v.MapIndex(key).Interface(), key.Interface()) != nil {
//				return
//			}
//		}
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		for i := 0; i < v.Len(); i++ {
//			if fn(v.Index(i).Interface(), i) != nil {
//				return
//			}
//		}
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}
//
//func Len(input interface{}) int {
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//		fallthrough
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		return v.Len()
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}
//
//func Cap(input interface{}) int {
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		return v.Cap()
//	default:
//		panic(fmt.Sprintf("input type not a slice or array: %v", v))
//	}
//}
//
//func MapToMap(input interface{}, fn Iteratee) map[interface{}]interface{} {
//	m := make(map[interface{}]interface{})
//	Each(input, func(val interface{}, key interface{}) interface{} {
//		m[key] = fn(val, key)
//		return nil
//	})
//	return m
//}
//
//func MapToSlice(input interface{}, fn Iteratee) []interface{} {
//	m := make([]interface{}, Len(input))
//	Each(input, func(val interface{}, key interface{}) interface{} {
//		m[key.(int)] = fn(val, key)
//		return nil
//	})
//	return m
//}
//
//func Filter(input interface{}, fn Iteratee) interface{} {
//	iType := reflect.TypeOf(input)
//	outVal := reflect.New(iType)
//	Each(input, func(val interface{}, key interface{}) interface{} {
//		if
//	})
//	return outVal.Interface()
//}
//
//func Range(start int, end int, step int) []int {
//	length := int(math.Ceil(float64(end-start) / float64(step)))
//	slc := make([]int, length)
//	for i := start; i < end; i += step {
//		slc = append(slc, i)
//	}
//	return slc
//}
//
//func Keys(input interface{}) []interface{} {
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		return Range()
//		keys := make([]interface{}, v.Len())
//
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//	if v.Kind() != reflect.Map {
//		panic(fmt.Sprintf("input type not a map: %v", v))
//	}
//	keys := make([]interface{}, 0, len(v.MapKeys()))
//	for _, k := range v.MapKeys() {
//		keys = append(keys, k.Interface())
//	}
//	return keys
//}
//
//func Vals(input interface{}) []interface{} {
//	v := reflect.ValueOf(input)
//	if v.Kind() != reflect.Map {
//		panic(fmt.Sprintf("input type not a map: %v", v))
//	}
//	vals := make([]interface{}, 0, len(v.MapKeys()))
//	for _, k := range v.MapKeys() {
//		vals = append(vals, v.MapIndex(k).Interface())
//	}
//	return vals
//}
//
//func Has(input interface{}, value interface{}) bool {
//	v := reflect.ValueOf(input)
//	if v.Kind() == reflect.Map {
//		input = Vals(input)
//	}
//	v = reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		for i := 0; i < v.Len(); i++ {
//			if reflect.DeepEqual(v.Index(i).Interface(), value) {
//				return true
//			}
//		}
//		return false
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}
//
//func ToSliceOfIface(input interface{}) []interface{} {
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//		return Vals(input)
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		slc := make([]interface{}, 0, v.Len())
//		for i := 0; i < v.Len(); i++ {
//			slc = append(slc, v.Index(i).Interface())
//		}
//		return slc
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}
//
//func ToSliceOfString(input interface{}) []string {
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//		return ToSliceOfString(Vals(input))
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		slc := make([]string, 0, v.Len())
//		for i := 0; i < v.Len(); i++ {
//			slc = append(slc, fmt.Sprintf("%v", v.Index(i).Interface()))
//		}
//		return slc
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}
//
//// WithoutKeys returns copy of the provided input (map, slice or array) without
//// elements with provided keys.
//func WithoutKeys(input interface{}, keys ...interface{}) interface{} {
//	kv := reflect.ValueOf(keys[0])
//	if kv.Kind() == reflect.Slice || kv.Kind() == reflect.Array {
//		keys = kv.Interface().([]interface{})
//	}
//
//	v := reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//		newMap := make(map[interface{}]interface{})
//		for _, key := range v.MapKeys() {
//			if !Has(keys, key) {
//				newMap[key] = v.MapIndex(key).Interface()
//			}
//		}
//		return newMap
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		newSlice := make([]interface{}, 0)
//		for i := 0; i < v.Len(); i++ {
//			if !Has(keys, i) {
//				newSlice = append(newSlice, v.Index(i).Interface())
//			}
//		}
//		return newSlice
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}
//
//// WithoutValues returns copy of the provided input (map, slice or array) without
//// elements with provided keys.
//func WithoutValues(input interface{}, values ...interface{}) interface{} {
//	var v reflect.Value
//
//	v = reflect.ValueOf(values[0])
//	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
//		values = v.Interface().([]interface{})
//	}
//
//	v = reflect.ValueOf(input)
//	switch v.Kind() {
//	case reflect.Map:
//		newMap := make(map[interface{}]interface{})
//		for _, key := range v.MapKeys() {
//			val := v.MapIndex(key).Interface()
//			if !Has(values, val) {
//				newMap[key] = val
//			}
//		}
//		return newMap
//	case reflect.Array:
//		fallthrough
//	case reflect.Slice:
//		newSlice := make([]interface{}, 0)
//		for i := 0; i < v.Len(); i++ {
//			val := v.Index(i).Interface()
//			if !Has(values, val) {
//				newSlice = append(newSlice, val)
//			}
//		}
//		return newSlice
//	default:
//		panic(fmt.Sprintf("input type not a slice, map or array: %v", v))
//	}
//}

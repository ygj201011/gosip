package collection

import (
	"reflect"
)

func isCollection(kind reflect.Kind) bool {
	for _, k := range collectionKinds {
		if kind == k {
			return true
		}
	}
	return false
}

// IsCollection checks if value is an array, slice, map or string.
func IsCollection(value Any) bool {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	return isCollection(valueType.Kind())
}

// IsFunction checks if value is a function.
func IsFunction(value Any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Func
}

// IsIteratee checks if value is a valid iteratee function.
func IsIteratee(value Any) bool {
	if !IsFunction(value) {
		return false
	}

	fnType := reflect.TypeOf(value)
	numIn := fnType.NumIn()
	numOut := fnType.NumOut()

	if numIn < 1 || numIn > 3 || numOut < 1 || numOut > 2 {
		return false
	}
	// check that second output param is of error type
	if numOut == 2 {
		a, b := fnType.Out(1), reflect.TypeOf((*error)(nil)).Elem()
		if !a.AssignableTo(b) {
			return false
		}
	}

	return true
}

// IsMapLikeIteratee checks if value is a valid map-like iteratee function.
func IsMapLikeIteratee(value Any) bool {
	if !IsIteratee(value) {
		return false
	}

	fnType := reflect.TypeOf(value)
	numIn := fnType.NumIn()

	if numIn > 2 {
		return false
	}

	return true
}

// IsReduceLikeIteratee checks if value is a valid map-like iteratee function.
func IsReduceLikeIteratee(value Any) bool {
	if !IsIteratee(value) {
		return false
	}

	fnType := reflect.TypeOf(value)
	numIn := fnType.NumIn()

	if numIn < 2 {
		return false
	}

	return true
}

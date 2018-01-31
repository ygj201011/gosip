package generic

import (
	"fmt"
	"reflect"
)

// IsPointer checks if value is an pointer.
func IsPointer(value Any) bool {
	return isPointer(reflect.ValueOf(value))
}

func isPointer(val reflect.Value) bool {
	return val.Kind() == reflect.Ptr
}

func isInterface(val reflect.Value) bool {
	return val.Kind() == reflect.Interface
}

// IsCollection checks if value is an array, slice, map or chan
// or pointer to array, slice, map or chan.
func IsCollection(value Any) bool {
	return isCollection(reflect.ValueOf(value))
}

func isCollection(val reflect.Value) bool {
	val = getActualValue(val)
	for _, k := range collectionKinds {
		if val.Kind() == k {
			return true
		}
	}
	return false
}

// IsCollectionPointer checks if value is an pointer to array, slice, map or chan.
func IsCollectionPointer(value Any) bool {
	return isCollectionPointer(reflect.ValueOf(value))
}

func isCollectionPointer(val reflect.Value) bool {
	if !isPointer(val) {
		return false
	}
	return isCollection(val)
}

// IsFunction checks if value is a function.
func IsFunction(value Any) bool {
	return isFunction(reflect.ValueOf(value))
}

func isFunction(val reflect.Value) bool {
	return val.Kind() == reflect.Func
}

// IsIteratee checks if value is a valid iteratee function.
func IsIteratee(value Any) bool {
	return isIteratee(reflect.ValueOf(value))
}

func isIteratee(val reflect.Value) bool {
	if !isFunction(val) {
		return false
	}

	typ := val.Type()
	numIn := typ.NumIn()
	numOut := typ.NumOut()

	if numIn < 1 || numIn > 3 || numOut < 1 || numOut > 2 {
		return false
	}
	// check that second output param is of error type
	if numOut == 2 {
		a, b := typ.Out(1), reflect.TypeOf((*error)(nil)).Elem()
		if !a.AssignableTo(b) {
			return false
		}
	}

	return true
}

// IsMapLikeIteratee checks if value is a valid map-like iteratee function.
func IsMapLikeIteratee(value Any) bool {
	return isMapLikeIteratee(reflect.ValueOf(value))
}

func isMapLikeIteratee(val reflect.Value) bool {
	if !isIteratee(val) {
		return false
	}

	typ := val.Type()
	numIn := typ.NumIn()
	if numIn > 2 {
		return false
	}

	return true
}

// IsReduceLikeIteratee checks if value is a valid map-like iteratee function.
func IsReduceLikeIteratee(value Any) bool {
	return isReduceLikeIteratee(reflect.ValueOf(value))
}

func isReduceLikeIteratee(val reflect.Value) bool {
	if !isIteratee(val) {
		return false
	}

	typ := val.Type()
	numIn := typ.NumIn()
	if numIn < 2 {
		return false
	}

	return true
}

func newCollectionValue(collectionType reflect.Type) (reflect.Value, error) {
	switch collectionType.Kind() {
	case reflect.Map:
		return reflect.MakeMap(collectionType), nil
	case reflect.Array:
		return reflect.New(collectionType).Elem(), nil
	case reflect.Slice:
		return reflect.MakeSlice(collectionType, 0, 0), nil
	case reflect.Chan:
		return reflect.MakeChan(collectionType, 0), nil
	default:
		return reflect.ValueOf(nil), NewUnexpectedTypeError("unexpected type of collection to create: expect %v, but got %v", collectionKinds, collectionType)
	}
}

func finalizeCollection(collectionValue *reflect.Value) {
	switch collectionValue.Kind() {
	case reflect.Chan:
		collectionValue.Close()
	}
}

type iteratee func(elemValue, keyValue reflect.Value) error

func iterate(collectionValue reflect.Value, iteratee iteratee) error {
	switch collectionValue.Kind() {
	case reflect.Map:
		for _, keyValue := range collectionValue.MapKeys() {
			elemValue := collectionValue.MapIndex(keyValue)
			if err := iteratee(elemValue, keyValue); err != nil {
				return err
			}
		}
		return nil
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < collectionValue.Len(); i++ {
			keyValue := reflect.ValueOf(i)
			elemValue := collectionValue.Index(i)
			if err := iteratee(elemValue, keyValue); err != nil {
				return err
			}
		}
		return nil
	case reflect.Chan:
		i := 0
		for {
			keyValue := reflect.ValueOf(i)
			elemValue, ok := collectionValue.Recv()
			if !ok {
				// TODO return closed chan error?
				return nil
			}
			if err := iteratee(elemValue, keyValue); err != nil {
				return err
			}
		}
		return nil
	default:
		return NewUnexpectedTypeError(MsgUnexpectedTypeOfInputCollection, collectionValue.Interface(), collectionKinds, collectionValue.Type())
	}
}

func appendToCollection(collectionValue *reflect.Value, elemValue, keyValue reflect.Value) error {
	switch collectionValue.Kind() {
	case reflect.Map:
		if collectionValue.IsNil() {
			return fmt.Errorf("collection of type %v has nil value", collectionValue.Type())
		}
		collectionValue.SetMapIndex(keyValue, elemValue)
		return nil
	case reflect.Array:
		key := keyValue.Interface().(int)
		if key >= collectionValue.Len() {
			return NewIndexOutOfRangeError(MsgIndexOutOfRangeOfResultCollectionType, key, collectionValue.Type())
		}
		if !collectionValue.Index(key).CanSet() {
			return fmt.Errorf("value at key %v of the collection of type %v isn't settable", key, collectionValue.Type())
		}
		collectionValue.Index(key).Set(elemValue)
		return nil
	case reflect.Slice:
		*collectionValue = reflect.Append(*collectionValue, elemValue)
		return nil
	case reflect.Chan:
		collectionValue.Send(elemValue)
		return nil
	default:
		return NewUnexpectedTypeError(MsgUnexpectedTypeOfInputCollection, collectionValue.Interface(), collectionKinds, collectionValue.Type())
	}
}

func verifyCollectionLengths(fromCollectionValue, toCollectionValue reflect.Value) error {
	switch toCollectionValue.Kind() {
	case reflect.Array:
		if toCollectionValue.Type().Len() < fromCollectionValue.Len() {
			return NewMismatchedError(MsgResultCollectionLenLessThanInputCollectionLen, toCollectionValue.Type().Len(), fromCollectionValue.Len())
		}
	}
	return nil
}

func verifyMapTypeAssignability(iterateeValue, collectionValue, resultValue reflect.Value) error {
	iterateeType := iterateeValue.Type()
	collectionType := collectionValue.Type()
	resultType := resultValue.Type()
	// check that result collection can store all elements from input collection
	if err := verifyCollectionLengths(collectionValue, resultValue); err != nil {
		return err
	}
	// check that input collection element type is assignable to the first iteratee argument
	if !collectionType.Elem().AssignableTo(iterateeType.In(0)) {
		return NewNotAssignableTypeError(MsgInputCollectionElemTypeNotAssignableToIterateeArgType, collectionType.Elem(), "1st", iterateeType.In(0))
	}
	// check input collection key type
	if collectionKeyType, err := getCollectionKeyType(collectionType); err != nil {
		return err
	} else {
		// if the second argument present in iteratee signature, i.e. key arg
		// check that collection key type is assignable to it
		if iterateeType.NumIn() == 2 {
			if !collectionKeyType.AssignableTo(iterateeType.In(1)) {
				return NewNotAssignableTypeError(MsgInputCollectionKeyTypeNotAssignableToIterateeArgType, collectionKeyType, "2nd", iterateeType.In(1))
			}
		}
		// check that collection key type is assignable to result collection key type
		// for cases like: slice -> map, map -> slice and etc.
		if resultCollectionKeyType, err := getCollectionKeyType(resultType); err != nil {
			return err
		} else {
			if !collectionKeyType.AssignableTo(resultCollectionKeyType) {
				return NewNotAssignableTypeError(MsgInputCollectionKeyTypeNotAssignableToResultCollectionKeyType, collectionKeyType, resultCollectionKeyType)
			}
		}
	}
	// check that first iteratee output type is assignable to result collection element type
	if !iterateeType.Out(0).AssignableTo(resultType.Elem()) {
		return NewNotAssignableTypeError(MsgIterateeOutputTypeNotAssignableToResultCollectionElemType, iterateeType.Out(0), resultType.Elem())
	}
	return nil
}

func verifyReduceTypeAssignability(iterateeValue, collectionValue, resultValue reflect.Value) error {
	iterateeType := iterateeValue.Type()
	collectionType := collectionValue.Type()
	resultType := resultValue.Type()

	// check that input collection element type is assignable to the second iteratee argument
	if !collectionType.Elem().AssignableTo(iterateeType.In(1)) {
		return NewNotAssignableTypeError(MsgInputCollectionElemTypeNotAssignableToIterateeArgType, collectionType.Elem(), "2nd", iterateeType.In(1))
	}
	// check input collection key type
	if collectionKeyType, err := getCollectionKeyType(collectionType); err != nil {
		return err
	} else {
		// if the third argument present in iteratee signature, i.e. key arg
		// check that collection key type is assignable to it
		if iterateeType.NumIn() == 3 {
			if !collectionKeyType.AssignableTo(iterateeType.In(2)) {
				return NewNotAssignableTypeError(MsgInputCollectionKeyTypeNotAssignableToIterateeArgType, collectionKeyType, "3rd", iterateeType.In(2))
			}
		}
	}
	// check iteratee output type is assignable to 1st input argument type
	if !iterateeType.Out(0).AssignableTo(iterateeType.In(0)) {
		return NewNotAssignableTypeError(MsgIterateeOutputTypeNotAssignableToIterateeArgType, iterateeType.Out(0), "1st", iterateeType.In(0))
	}
	// check iteratee output type is assignable to result type
	if !iterateeType.Out(0).AssignableTo(resultType) {
		return NewNotAssignableTypeError(MsgIterateeOutputTypeNotAssignableToResultType, iterateeType.Out(0), resultType)
	}
	return nil
}

func getCollectionKeyType(collectionType reflect.Type) (reflect.Type, error) {
	switch collectionType.Kind() {
	case reflect.Map:
		return collectionType.Key(), nil
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Chan:
		return reflect.TypeOf(int(0)), nil
	default:
		return nil, NewUnexpectedTypeError(MsgUnexpectedTypeOfInputCollection, collectionType, collectionKinds, collectionType)
	}
}

func getActualValue(val reflect.Value) reflect.Value {
	if isPointer(val) {
		val = val.Elem()
	}
	if isInterface(val) {
		val = val.Elem()
	}
	return val
}

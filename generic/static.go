package generic

import (
	"reflect"
)

// Map creates new collection by running each element in the input collection
// through the map-like iteratee. Result will be written to value of the resultPtr argument, that
// should be pointer to collection-like value.
func Map(collection Collection, iteratee Iteratee, resultPtr Pointer) error {
	// collection reflection
	collectionValue := reflect.ValueOf(collection)
	// iteratee reflection
	iterateeValue := reflect.ValueOf(iteratee)
	// result reflection
	resultPtrValue := reflect.ValueOf(resultPtr)
	// check input types
	if !isCollection(collectionValue) {
		return NewUnexpectedTypeError(MsgUnexpectedTypeOfInputCollection, collection, collectionKinds, collectionValue.Type())
	}
	if !isMapLikeIteratee(iterateeValue) {
		return NewUnexpectedTypeError(MsgUnexpectedTypeOfMapLikeIteratee, iterateeValue.Interface(), PseudoMapLikeIteratee, iterateeValue.Type())
	}
	if !isCollectionPointer(resultPtrValue) {
		return NewUnexpectedTypeError(MsgUnexpectedTypeOfResultPointer, resultPtr, collectionPtrKinds, resultPtrValue.Type())
	}
	iterateeType := iterateeValue.Type()
	// get actual value of input / output collections
	collectionValue = getActualValue(collectionValue)
	resultValue := getActualValue(resultPtrValue)
	// and verify that all related arguments have consistent types
	if err := verifyMapTypeAssignability(iterateeValue, collectionValue, resultValue); err != nil {
		return err
	}
	// create new collection of result type
	var err error
	newCollectionValue, err := newCollectionValue(resultValue.Type())
	if err != nil {
		return err
	}

	err = iterate(collectionValue, func(elemValue reflect.Value, keyValue reflect.Value) error {
		var iterateeResultValues []reflect.Value
		// invoke external iteratee
		switch iterateeType.NumIn() {
		case 1:
			iterateeResultValues = iterateeValue.Call([]reflect.Value{elemValue})
		case 2:
			iterateeResultValues = iterateeValue.Call([]reflect.Value{elemValue, keyValue})
		}

		if iterateeType.NumOut() == 2 {
			if err := iterateeResultValues[1]; !err.IsNil() {
				return err.Interface().(error)
			}
		}

		return appendToCollection(&newCollectionValue, iterateeResultValues[0], keyValue)
	})

	finalizeCollection(&newCollectionValue)
	resultPtrValue.Elem().Set(newCollectionValue)

	return err
}

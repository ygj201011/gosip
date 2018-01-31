package generic

import (
	"errors"
	"fmt"
)

// Package common error messages.
const (
	MsgUnexpectedTypeOfInputCollection                              = "unexpected type of the input collection %v: expect %v, but got %v"
	MsgUnexpectedTypeOfMapLikeIteratee                              = "unexpected type of the map-like iteratee %v: expect %v, but got %v"
	MsgUnexpectedTypeOfResultPointer                                = "unexpected type of the result pointer %v: expect %v, but got %v"
	MsgInputCollectionElemTypeNotAssignableToIterateeArgType        = "an element type %v of the input collection isn't assignable to the %v iteratee argument type %v"
	MsgInputCollectionKeyTypeNotAssignableToIterateeArgType         = "a key type %v of the input collection isn't assignable to the %v iteratee argument type %v"
	MsgInputCollectionKeyTypeNotAssignableToResultCollectionKeyType = "a key type %v of the input collection isn't assignable to the result collection key %v"
	MsgIterateeOutputTypeNotAssignableToResultCollectionElemType    = "the 1st output type %v of the iteratee isn't assignable to the result collection element type %v"
	MsgIterateeOutputTypeNotAssignableToResultType                  = "the 1st output type %v of the iteratee isn't assignable to the result type %v"
	MsgIterateeOutputTypeNotAssignableToIterateeArgType             = "the 1st output type %v of the iteratee isn't assignable to the %v iteratee argument type %v"
	MsgIndexOutOfRangeOfResultCollectionType                        = "index %v is out of range of the result collection with type %v"
	MsgResultCollectionLenLessThanInputCollectionLen                = "length %v of the result collection less than length %v of the input collection"
)

// ErrBreak is a special error that can be used to break collection iteration.
var ErrBreak = errors.New("breaked out")

// UnexpectedTypeError is used when method, that expects collection-like type, receives something another.
type UnexpectedTypeError string

// NewUnexpectedTypeError creates new UnexpectedTypeError with formatted message.
func NewUnexpectedTypeError(msg string, args ...interface{}) UnexpectedTypeError {
	return UnexpectedTypeError(fmt.Sprintf(msg, args...))
}
func (err UnexpectedTypeError) Error() string { return string(err) }

// NotAssignableTypeError is used when one type isn't assignable to another.
type NotAssignableTypeError string

// NewNotAssignableTypeError creates new NotAssignableTypeError with formatted message.
func NewNotAssignableTypeError(msg string, args ...interface{}) NotAssignableTypeError {
	return NotAssignableTypeError(fmt.Sprintf(msg, args...))
}
func (err NotAssignableTypeError) Error() string { return string(err) }

// IndexOutOfRangeError is used when provided collection key is out of maximum range.
type IndexOutOfRangeError string

// NewIndexOutOfRangeError creates new IndexOutOfRangeError with formatted message.
func NewIndexOutOfRangeError(msg string, args ...interface{}) IndexOutOfRangeError {
	return IndexOutOfRangeError(fmt.Sprintf(msg, args...))
}
func (err IndexOutOfRangeError) Error() string { return string(err) }

// MismatchedError is used when some value is mismatched with another value.
type MismatchedError string

// NewMismatchedError creates new MismatchedError with formatted message.
func NewMismatchedError(msg string, args ...interface{}) MismatchedError {
	return MismatchedError(fmt.Sprintf(msg, args...))
}
func (err MismatchedError) Error() string { return string(err) }

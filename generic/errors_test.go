package generic_test

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ghettovoice/gosip/generic"
)

var _ = Describe("Errors", func() {
	Context("NewUnexpectedTypeError", func() {
		It("should return new formatted UnexpectedTypeError", func() {
			args := []interface{}{reflect.TypeOf([]int{})}
			act := NewUnexpectedTypeError("unexpected type %v received", args...)
			exp := fmt.Sprintf("unexpected type %v received", args...)

			Expect(act).To(MatchError(exp))
		})
	})

	Context("NewNotAssignableTypeError", func() {
		It("should return new formatted NotAssignableTypeError", func() {
			args := []interface{}{reflect.TypeOf([]int{}), reflect.TypeOf(0)}
			act := NewNotAssignableTypeError("type %v isn't assignable to type %v", args...)
			exp := fmt.Sprintf("type %v isn't assignable to type %v", args...)

			Expect(act).To(MatchError(exp))
		})
	})

	Context("NewIndexOutOfRangeError", func() {
		It("should return new formatted IndexOutOfRangeError", func() {
			args := []interface{}{5, reflect.TypeOf([3]int{1, 2, 3})}
			act := NewIndexOutOfRangeError("key %v is out of range of type %v", args...)
			exp := fmt.Sprintf("key %v is out of range of type %v", args...)

			Expect(act).To(MatchError(exp))
		})
	})

	Context("NewMismatchedError", func() {
		It("should return new formatted MismatchedError", func() {
			args := []interface{}{2, 3}
			act := NewMismatchedError("length %v is less than length %v", args...)
			exp := fmt.Sprintf("length %v is less than length %v", args...)

			Expect(act).To(MatchError(exp))
		})
	})
})

package generic_test

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/ghettovoice/gosip/generic"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Static API", func() {
	var err error

	Describe("Map()", func() {
		Context("given unexpected arguments", func() {
			Context("unexpected input collection type", func() {
				st := struct{ id int }{id: 123}
				inputs := []Any{123, "qwerty", true, st, &st}
				entries := make([]TableEntry, 0)

				for _, in := range inputs {
					t := reflect.TypeOf(in)
					entries = append(entries,
						Entry(t.String(), in, NewUnexpectedTypeError(MsgUnexpectedTypeOfInputCollection, in, CollectionKinds(), t)),
					)
				}

				DescribeTable("type validation",
					func(in Any, expected UnexpectedTypeError) {
						var res []int
						iter := func(val int) int { return val }

						err = Map(in, iter, &res)
						Expect(err).To(MatchError(expected))
						Expect(res).To(HaveLen(0))
					},
					entries...,
				)
			})

			Context("unexpected iteratee type", func() {
				iteratees := []Any{
					123,
					"not a function",
					func() {},
					func(val int) {},
					func(val string, key int) {},
					func(val string) (string, int) { return "", 0 },
					func(val string) (string, error, int) { return "", nil, 0 },
					func(val string, key int) (string, error, int) { return "", nil, 0 },
				}
				entries := make([]TableEntry, 0)

				for _, iter := range iteratees {
					t := reflect.TypeOf(iter)
					entries = append(entries,
						Entry(t.String(), iter, NewUnexpectedTypeError(MsgUnexpectedTypeOfMapLikeIteratee, iter, PseudoMapLikeIteratee, t)),
					)
				}

				DescribeTable("type validation",
					func(iter Any, expected UnexpectedTypeError) {
						var res []string
						in := []string{"q", "w", "e"}

						err = Map(in, iter, &res)
						Expect(err).To(MatchError(expected))
						Expect(res).To(HaveLen(0))
					},
					entries...,
				)
			})

			Context("unexpected result pointer type", func() {
				str := "qwerty"
				num := 123
				fn := func() {}
				bl := true
				struc := struct{ id int }{id: 123}
				results := []Any{str, &str, num, &num, fn, &fn, bl, &bl, struc, &struc}
				entries := make([]TableEntry, 0)

				for _, res := range results {
					t := reflect.TypeOf(res)
					entries = append(entries,
						Entry(t.String(), res, NewUnexpectedTypeError(MsgUnexpectedTypeOfResultPointer, res, CollectionPtrKinds(), t)),
					)
				}

				DescribeTable("type validation",
					func(res Any, expected UnexpectedTypeError) {
						in := []string{"q", "w", "e"}
						iter := func(val string, key int) string { return val }

						err = Map(in, iter, res)
						Expect(err).To(MatchError(expected))
					},
					entries...,
				)
			})

			Context("valid but inconsistent arguments", func() {
				var in Any
				var iter Iteratee
				var err error

				Context("collection element type isn't assignable to the 1st iteratee argument type", func() {
					var res []int

					BeforeEach(func() {
						in = []string{"q", "w", "e"}
						iter = func(val int) int { return val }
						res = make([]int, 0)

						err = Map(in, iter, &res)
					})

					It("should return NotAssignableError", func() {
						Expect(err).To(MatchError(NewNotAssignableTypeError(
							MsgInputCollectionElemTypeNotAssignableToIterateeArgType,
							reflect.TypeOf(in).Elem(),
							"1st",
							reflect.TypeOf(iter).In(0),
						)))
					})

					It("should not change result", func() {
						Expect(res).To(HaveLen(0))
					})
				})

				Context("collection key type isn't assignable to the 2nd iteratee argument", func() {
					var res []string

					BeforeEach(func() {
						in = []string{"q", "w", "e"}
						iter = func(val string, key string) string { return val }
						res = make([]string, 0)

						err = Map(in, iter, &res)
					})

					It("should return NotAssignableError", func() {
						Expect(err).To(MatchError(NewNotAssignableTypeError(
							MsgInputCollectionKeyTypeNotAssignableToIterateeArgType,
							reflect.TypeOf(0),
							"2nd",
							reflect.TypeOf(iter).In(1),
						)))
					})

					It("should not change result", func() {
						Expect(res).To(HaveLen(0))
					})
				})

				Context("input collection key type isn't assignable to result collection key type", func() {
					var res []string

					BeforeEach(func() {
						in = map[string]string{
							"one":   "q",
							"two":   "w",
							"three": "e",
						}
						iter = func(val string, key string) string { return val }
						res = make([]string, 0)

						err = Map(in, iter, &res)
					})

					It("should return NotAssignableTypeError", func() {
						Expect(err).To(MatchError(NewNotAssignableTypeError(
							MsgInputCollectionKeyTypeNotAssignableToResultCollectionKeyType,
							reflect.TypeOf(in).Key(),
							reflect.TypeOf(0),
						)))
					})

					It("should not change result", func() {
						Expect(res).To(HaveLen(0))
					})
				})

				Context("result collection can't hold all input elements", func() {
					Context("result is an array of fixed length less than input collection length", func() {
						var res [2]string

						BeforeEach(func() {
							in = []string{"q", "w", "e"}
							iter = func(val string) string { return val }
							res = *new([2]string)
							err = Map(in, iter, &res)
						})

						It("should return MismatchedError", func() {
							Expect(err).To(MatchError(NewMismatchedError(MsgResultCollectionLenLessThanInputCollectionLen, 2, 3)))
						})

						It("should not change result", func() {
							for _, val := range res {
								Expect(val).To(BeEmpty())
							}
						})
					})
				})

				Context("iteratee 1st output type isn't assignable to result collection element type", func() {
					var res []string
					BeforeEach(func() {
						in = []int{1, 2, 3}
						iter = func(val int, key int) int { return val }

						err = Map(in, iter, &res)
					})

					It("should return NotAssignableTypeError", func() {
						Expect(err).To(MatchError(NewNotAssignableTypeError(
							MsgIterateeOutputTypeNotAssignableToResultCollectionElemType,
							reflect.TypeOf(iter).Out(0),
							reflect.TypeOf(res).Elem(),
						)))
					})

					It("should not change result", func() {
						Expect(res).To(HaveLen(0))
					})
				})
			})
		})

		Context("given valid arguments", func() {
			Context("given iteratee defined with 2 outputs", func() {
				var res []string
				var iter func(val string, key int) (string, error)
				in := []string{"h", "e", "l", "l", "o"}

				BeforeEach(func() {
					res = make([]string, 0)
				})

				JustBeforeEach(func() {
					err = Map(in, iter, &res)
				})

				Context("applied without error", func() {
					BeforeEach(func() {
						iter = func(val string, key int) (string, error) {
							return val + val, nil
						}
					})

					It("should not return error", func() {
						Expect(err).ToNot(HaveOccurred())
					})

					It("should write mapped collection", func() {
						Expect(res).To(Equal([]string{"hh", "ee", "ll", "ll", "oo"}))
					})
				})

				Context("applied with error", func() {
					breakErr := errors.New("failed")

					BeforeEach(func() {
						iter = func(val string, key int) (string, error) {
							if key == 2 {
								return "", breakErr
							}
							return val + val, nil
						}
					})

					It("should return error from failed iteration", func() {
						Expect(err).To(MatchError(breakErr))
					})

					It("should write result of successful iterations", func() {
						Expect(res).To(Equal([]string{"hh", "ee"}))
					})
				})
			})

			DescribeTable("normal work",
				func(in Collection, iter Iteratee, res Collection, expected Collection) {
					err = Map(in, iter, &res)

					Expect(err).ToNot(HaveOccurred())
					Expect(res).To(Equal(expected))
				},
				Entry("array -> array",
					[...]int{1, 2, 3},
					func(val int, key int) string { return fmt.Sprintf("%d -> %d", key, val) },
					[3]string{},
					[...]string{"0 -> 1", "1 -> 2", "2 -> 3"},
				),
			)
			//DescribeTable("array -> array",
			//	func(in Collection, iter Iteratee, res Collection, expected Collection) {
			//		var res [len]
			//		Expect(Map())
			//	},
			//)
		})

		//Context("when input is an array", func() {
		//	in := [...]int{1, 2, 3}
		//	iter := func(val int, key int) string {
		//		return fmt.Sprintf("val %v + key %v = %v", val, key, val+key)
		//	}
		//
		//	It("should produce new array with the same length as input collection", func() {
		//		var out [3]string
		//		exp := [...]string{
		//			"val 1 + key 0 = 1",
		//			"val 2 + key 1 = 3",
		//			"val 3 + key 2 = 5",
		//		}
		//
		//		Expect(Map(in, iter, &out)).ToNot(HaveOccurred())
		//		Expect(out).To(HaveLen(len(in)))
		//		Expect(out).To(Equal(exp))
		//	})
		//
		//	It("should produce new slice with the same length as input collection", func() {
		//		out := make([]string, 0)
		//		exp := []string{
		//			"val 1 + key 0 = 1",
		//			"val 2 + key 1 = 3",
		//			"val 3 + key 2 = 5",
		//		}
		//
		//		Expect(Map(in, iter, &out)).ToNot(HaveOccurred())
		//		Expect(out).To(HaveLen(len(in)))
		//		Expect(out).To(Equal(exp))
		//	})
		//
		//	It("should produce new map with same length and transformed elements", func() {
		//		out := make(map[int]string)
		//		exp := map[int]string{
		//			0: "val 1 + key 0 = 1",
		//			1: "val 2 + key 1 = 3",
		//			2: "val 3 + key 2 = 5",
		//		}
		//
		//		Expect(Map(in, iter, &out)).ToNot(HaveOccurred())
		//		Expect(out).To(HaveLen(len(in)))
		//		Expect(out).To(Equal(exp))
		//	})
		//
		//	It("should produce new chan that emits transformed elements", func(done Done) {
		//		defer close(done)
		//
		//		out := make(chan string)
		//		exp := []string{
		//			"val 1 + key 0 = 1",
		//			"val 2 + key 1 = 3",
		//			"val 3 + key 2 = 5",
		//		}
		//
		//		go func() {
		//			Expect(Map(in, iter, &out)).ToNot(HaveOccurred())
		//		}()
		//
		//		res := make([]string, 0)
		//		for val := range out {
		//			res = append(res, val)
		//		}
		//
		//		Expect(res).To(HaveLen(len(in)))
		//		Expect(res).To(Equal(exp))
		//	})
		//})
	})
})

func ExampleMap() {
	strs := make([]string, 0)
	input := []int{1, 2, 3}
	iteratee := func(val int, key int) string { return fmt.Sprintf("%d = %d", key, val*2) }
	err := Map(input, iteratee, &strs)

	fmt.Println(err)
	fmt.Printf("%v\n", strs)
	// Output:
	// <nil>
	// [0 = 2 1 = 4 2 = 6]
}

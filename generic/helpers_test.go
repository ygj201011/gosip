package generic_test

import (
	. "github.com/ghettovoice/gosip/generic"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {
	Context("IsPointer()", func() {
		var s string
		var n int
		var slc []float32

		It("should return true for pointer values", func() {
			Expect(IsPointer(&s)).To(BeTrue())
			Expect(IsPointer(&n)).To(BeTrue())
			Expect(IsPointer(&slc)).To(BeTrue())
		})

		It("should return false for non pointer values", func() {
			Expect(IsPointer(s)).To(BeFalse())
			Expect(IsPointer(n)).To(BeFalse())
			Expect(IsPointer(slc)).To(BeFalse())
		})
	})

	Context("IsCollection()", func() {
		It("should return true for array and pointer to array values", func() {
			in := [...]int{1, 2, 3}
			Expect(IsCollection(in)).To(BeTrue())
			Expect(IsCollection(&in)).To(BeTrue())
		})

		It("should return true for slice and pointer to slice values", func() {
			in := []string{"q", "w", "e"}
			Expect(IsCollection(in)).To(BeTrue())
			Expect(IsCollection(&in)).To(BeTrue())
		})

		It("should return true for map or pointer to map values", func() {
			in := map[int]string{1: "q", 2: "w", 3: "e"}
			Expect(IsCollection(in)).To(BeTrue())
			Expect(IsCollection(&in)).To(BeTrue())
		})

		It("should return true for chan or pointer to chan values", func() {
			in := make(chan int)
			Expect(IsCollection(in)).To(BeTrue())
			Expect(IsCollection(&in)).To(BeTrue())
		})

		It("should return false for non collection values", func() {
			in := "qwerty"
			Expect(IsCollection(in)).To(BeFalse())
			Expect(IsCollection(&in)).To(BeFalse())
			Expect(IsCollection(123)).To(BeFalse())
			Expect(IsCollection(&struct{ id int }{123})).To(BeFalse())
			Expect(IsCollection(func() {})).To(BeFalse())
			Expect(IsCollection(true)).To(BeFalse())
		})
	})

	Context("IsCollectionPointer()", func() {
		It("should return true for pointer to array values", func() {
			in := [...]int{1, 2, 3}
			Expect(IsCollectionPointer(&in)).To(BeTrue())
		})

		It("should return true for pointer to slice values", func() {
			in := []string{"q", "w", "e"}
			Expect(IsCollectionPointer(&in)).To(BeTrue())
		})

		It("should return true for chan or pointer to chan values", func() {
			in := make(chan int)
			Expect(IsCollectionPointer(&in)).To(BeTrue())
		})

		It("should return false for non pointer to collection values", func() {
			arr := [...]int{1, 2, 3}
			slc := []string{"q", "w", "e"}
			ch := make(chan int)
			str := "qwerty"
			Expect(IsCollectionPointer(arr)).To(BeFalse())
			Expect(IsCollectionPointer(slc)).To(BeFalse())
			Expect(IsCollectionPointer(ch)).To(BeFalse())
			Expect(IsCollectionPointer(&str)).To(BeFalse())
			Expect(IsCollectionPointer(true)).To(BeFalse())
		})
	})

	Context("IsFunction()", func() {
		It("should return true for function values", func() {
			Expect(IsFunction(func() {})).To(BeTrue())
		})

		It("should return false for non function values", func() {
			Expect(IsFunction(false)).To(BeFalse())
			Expect(IsFunction(1)).To(BeFalse())
			Expect(IsFunction("qwerty")).To(BeFalse())
		})
	})

	Context("IsIteratee()", func() {
		It("should return true for iteratee functions", func() {
			Expect(IsIteratee(func(val int) int { return 0 })).To(BeTrue())
			Expect(IsIteratee(func(val int) error { return nil })).To(BeTrue())
			Expect(IsIteratee(func(val int) (int, error) { return 0, nil })).To(BeTrue())
			Expect(IsIteratee(func(val string, key int) string { return "" })).To(BeTrue())
			Expect(IsIteratee(func(val string, key int) error { return nil })).To(BeTrue())
			Expect(IsIteratee(func(val string, key int) (string, error) { return "", nil })).To(BeTrue())
			Expect(IsIteratee(func(all []string, val string, key int) []string { return []string{} })).To(BeTrue())
			Expect(IsIteratee(func(all []string, val string, key int) ([]string, error) { return []string{}, nil })).To(BeTrue())
		})

		It("should return false for non iteratee functions", func() {
			Expect(IsIteratee(func() {})).To(BeFalse())
			Expect(IsIteratee(func(val int) {})).To(BeFalse())
			Expect(IsIteratee(func(val int) (string, int) { return "", 0 })).To(BeFalse())
			Expect(IsIteratee(func(val int) (string, error, int) { return "", nil, 0 })).To(BeFalse())
			Expect(IsIteratee(func(val int, key string) {})).To(BeFalse())
			Expect(IsIteratee(func(val int, key string) (string, string) { return "", "" })).To(BeFalse())
			Expect(IsIteratee(func(val int, key string) (string, string, error) { return "", "", nil })).To(BeFalse())
			Expect(IsIteratee(func(val int, key, arg, arg2 string) (string, error) { return "", nil })).To(BeFalse())
		})
	})

	Context("IsMapLikeIteratee()", func() {
		It("should return true for map-like iteratee functions", func() {
			Expect(IsMapLikeIteratee(func(val string) string { return "" })).To(BeTrue())
			Expect(IsMapLikeIteratee(func(val string) (string, error) { return "", nil })).To(BeTrue())
			Expect(IsMapLikeIteratee(func(val, key string) string { return "" })).To(BeTrue())
			Expect(IsMapLikeIteratee(func(val, key string) (string, error) { return "", nil })).To(BeTrue())
		})

		It("should return false for non-map-like iteratee functions", func() {
			Expect(IsMapLikeIteratee(func() {})).To(BeFalse())
			Expect(IsMapLikeIteratee(func(val string) {})).To(BeFalse())
			Expect(IsMapLikeIteratee(func(val string) (string, string) { return "", "" })).To(BeFalse())
			Expect(IsMapLikeIteratee(func(val string) (string, error, string) { return "", nil, "" })).To(BeFalse())
			Expect(IsMapLikeIteratee(func(val, key string) {})).To(BeFalse())
			Expect(IsMapLikeIteratee(func(val, key string) (string, string) { return "", "" })).To(BeFalse())
			Expect(IsMapLikeIteratee(func(val, key string) (string, string, error) { return "", "", nil })).To(BeFalse())
			Expect(IsMapLikeIteratee(func(all []string, val, key string) string { return "" })).To(BeFalse())
		})
	})

	Context("IsReduceLikeIteratee()", func() {
		It("should return true for reduce-like iteratee functions", func() {
			Expect(IsReduceLikeIteratee(func(res, val string) string { return "" })).To(BeTrue())
			Expect(IsReduceLikeIteratee(func(res, val string) (string, error) { return "", nil })).To(BeTrue())
			Expect(IsReduceLikeIteratee(func(res, val, key string) string { return "" })).To(BeTrue())
			Expect(IsReduceLikeIteratee(func(res, val, key string) (string, error) { return "", nil })).To(BeTrue())
			Expect(IsReduceLikeIteratee(func(res []string, val, key string) ([]string, error) { return []string{}, nil })).To(BeTrue())
		})

		It("should return false for non-reduce-like iteratee functions", func() {
			Expect(IsReduceLikeIteratee(func() {})).To(BeFalse())
			Expect(IsReduceLikeIteratee(func() int { return 0 })).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(val int) {})).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(val int) int { return 0 })).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res, val int) {})).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res, val int) (int, int) { return 0, 1 })).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res, val int) (int, int) { return 0, 1 })).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res, val int) (int, error, string) { return 0, nil, "" })).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res []int, val, key int) {})).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res []int, val, key int) ([]int, int) { return []int{}, 0 })).To(BeFalse())
			Expect(IsReduceLikeIteratee(func(res []int, val, key int) ([]int, error, string) { return []int{}, nil, "" })).To(BeFalse())
		})
	})
})

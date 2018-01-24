package collection_test

import (
	. "github.com/ghettovoice/gosip/collection"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {
	Context("IsCollection()", func() {
		It("should return true for array value", func() {
			input := [...]int{1, 2, 3}
			Expect(IsCollection(input)).To(BeTrue())
			Expect(IsCollection(&input)).To(BeTrue())
		})

		It("should return true for slice value", func() {
			input := []string{"q", "w", "e"}
			Expect(IsCollection(input)).To(BeTrue())
			Expect(IsCollection(&input)).To(BeTrue())
		})

		It("should return true for map value", func() {
			input := map[int]string{1: "q", 2: "w", 3: "e"}
			Expect(IsCollection(input)).To(BeTrue())
			Expect(IsCollection(&input)).To(BeTrue())
		})

		It("should return true for string value", func() {
			input := "qwe"
			Expect(IsCollection(input)).To(BeTrue())
			Expect(IsCollection(&input)).To(BeTrue())
		})

		It("should return false for non-collection value", func() {
			Expect(IsCollection(123)).To(BeFalse())
			Expect(IsCollection(&struct{ id int }{123})).To(BeFalse())
			Expect(IsCollection(func() {})).To(BeFalse())
			Expect(IsCollection(true)).To(BeFalse())
		})
	})

	Context("IsFunction()", func() {
		It("should return true for function value", func() {
			Expect(IsFunction(func() {})).To(BeTrue())
		})

		It("should return false for non-function value", func() {
			Expect(IsFunction(false)).To(BeFalse())
			Expect(IsFunction(1)).To(BeFalse())
			Expect(IsFunction("qwerty")).To(BeFalse())
		})
	})

	Context("IsIteratee()", func() {
		It("should return true for iteratee function", func() {
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

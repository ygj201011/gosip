package collection_test

import (
	"github.com/ghettovoice/gosip/collection"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Static", func() {
	PContext("Map()", func() {
		It("should fail on invalid input arguments", func() {
			var res []int
			Expect(collection.Map(123, func() {}, &res)).To(HaveOccurred())
		})

		Context("of the array", func() {
			It("should return copy of the array", func() {

			})
		})
		It("should apply iteratee to each element of the collection", func() {
			var res []string
			Expect(collection.Map([...]string{"q", "w", "e", "r", "t"}, func(val string) string { return val + val }, res)).ToNot(HaveOccurred())
			Expect(res).To(Equal([...]string{"qq", "ww", "ee", "rr", "tt"}))
		})
	})

	//Context("Has()", func() {
	//	It("should panic when input is not a slice, map or array", func() {
	//		Expect(func() {
	//			collutil.Has("qwerty", "q")
	//		}).To(Panic())
	//	})
	//
	//	Context("called on slice", func() {
	//		It("should return true for existent value", func() {
	//			for _, elem := range slc {
	//				Expect(collutil.Has(slc, elem)).To(BeTrue())
	//			}
	//		})
	//		It("should return false for non-existent value", func() {
	//			Expect(collutil.Has(slc, "bullshit")).ToNot(BeTrue())
	//		})
	//	})
	//
	//	Context("called on array", func() {
	//		It("should return true for existent value", func() {
	//			for _, elem := range arr {
	//				Expect(collutil.Has(arr, elem)).To(BeTrue())
	//			}
	//		})
	//		It("should return false for non-existent value", func() {
	//			Expect(collutil.Has(arr, "bullshit")).ToNot(BeTrue())
	//		})
	//	})
	//
	//	Context("called on map", func() {
	//		It("should return true for existent value", func() {
	//			for _, elem := range ma {
	//				Expect(collutil.Has(ma, elem)).To(BeTrue())
	//			}
	//		})
	//		It("should return false for non-existent value", func() {
	//			Expect(collutil.Has(ma, "bullshit")).ToNot(BeTrue())
	//		})
	//	})
	//})
	//
	//Context("WithoutKeys()", func() {
	//	It("should panic when input is not a slice, map or array", func() {
	//		Expect(func() {
	//			collutil.WithoutKeys("qwerty", "q")
	//		}).To(Panic())
	//	})
	//
	//	Context("called on slice", func() {
	//
	//	})
	//})
})

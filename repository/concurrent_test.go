package repository_test

import (
	"fmt"

	. "github.com/ghettovoice/gosip/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concurrent repository", func() {
	var repo Repository

	input := behaviorInput{}
	items1 := map[string]interface{}{
		"one":   "first",
		"two":   "second",
		"three": "third",
	}
	items2 := map[string]interface{}{
		"one":   "first",
		"three": "third",
	}
	items3 := map[string]interface{}{
		"two":   "second",
		"three": "third",
	}

	BeforeEach(func() {
		repo = NewConcurrent()
		input.repo = repo
	})
	AfterEach(func() {
		repo = nil
		input = behaviorInput{}
	})

	Context("when just initialized", func() {
		BeEmptyBehavior(&input)
	})

	Context(fmt.Sprintf("when items added: %v", items1), func() {
		BeforeEach(func() {
			input.expectedItems = items1
			input.key = "three"
			input.expectedItem = items1[input.key]
			for k, v := range items1 {
				repo.Put(k, v)
			}
		})

		HaveLengthBehavior(&input)
		HaveItemsBehavior(&input)
		ItemFoundBehaviour(&input)

		Context("drop item with key 'two'", func() {
			BeforeEach(func() {
				input.expectedItems = items2
				input.key = "two"
				repo.Drop("two")
			})

			HaveLengthBehavior(&input)
			HaveItemsBehavior(&input)
			ItemNotFoundBehaviour(&input)
		})

		Context("pop item with key 'one'", func() {
			var it interface{}
			var ok bool

			BeforeEach(func() {
				input.expectedItems = items3
				input.key = "one"
				it, ok = repo.Pop("one")
				Expect(ok).To(BeTrue())
				Expect(it).ToNot(BeNil())
			})

			HaveLengthBehavior(&input)
			HaveItemsBehavior(&input)
			ItemNotFoundBehaviour(&input)

			It("should return item 'first'", func() {
				Expect(it).To(Equal("first"))
			})
		})

		Context("clear repository", func() {
			BeforeEach(func() {
				repo.Clear()
			})

			BeEmptyBehavior(&input)
		})
	})
})

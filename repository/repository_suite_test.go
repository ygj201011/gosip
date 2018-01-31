package repository_test

import (
	"fmt"
	"testing"

	"github.com/ghettovoice/gosip/generic"
	. "github.com/ghettovoice/gosip/repository"
	. "github.com/ghettovoice/gosip/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
	InitTestSuite(t, "Repository Suite")
}

type behaviorInput struct {
	repo           Repository
	key            string
	expectedItem   interface{}
	expectedItems  map[string]interface{}
	expectedLength int
	expectedKeys   []string
}

func BeEmptyBehavior(input *behaviorInput) {
	It("should be empty", func() {
		ExpectWithOffset(1, input.repo.Len()).To(Equal(0))
		ExpectWithOffset(1, input.repo.Keys()).To(HaveLen(0))
	})
}

func HaveLengthBehavior(input *behaviorInput) {
	It(fmt.Sprintf("should have length = %d", input.expectedLength), func() {
		ExpectWithOffset(1, input.repo.Len()).To(Equal(input.expectedLength))
		ExpectWithOffset(1, input.repo.Keys()).To(HaveLen(input.expectedLength))
	})
}

func HaveItemsBehavior(input *behaviorInput) {
	It(fmt.Sprintf("should have items: %+v", input.expectedItems), func() {
		ExpectWithOffset(1, input.repo.Items()).To(ConsistOf(generic.Vals(input.expectedItems)))
		ExpectWithOffset(1, input.repo.All()).To(ConsistOf(generic.Vals(input.expectedItems)))
	})
	It(fmt.Sprintf("should have keys: %v", generic.Keys(input.expectedItems)), func() {
		ExpectWithOffset(1, input.repo.Keys()).To(ConsistOf(generic.Keys(input.expectedItems)))
	})
}

func ItemFoundBehaviour(input *behaviorInput) {
	It(fmt.Sprintf("should find item by key %s", input.key), func() {
		ExpectWithOffset(1, input.repo.Has(input.key)).To(BeTrue())

		item, ok := input.repo.Get(input.key)
		ExpectWithOffset(1, ok).To(BeTrue())
		ExpectWithOffset(1, item).To(Equal(input.expectedItem))
	})
}

func ItemNotFoundBehaviour(input *behaviorInput) {
	It(fmt.Sprintf("should not find item by key %s", input.key), func() {
		ExpectWithOffset(1, input.repo.Has(input.key)).ToNot(BeTrue())

		item, ok := input.repo.Get(input.key)
		ExpectWithOffset(1, ok).To(BeFalse())
		ExpectWithOffset(1, item).To(BeNil())
	})
}

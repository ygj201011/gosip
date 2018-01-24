package repository_test

import (
	"fmt"
	"testing"

	"github.com/ghettovoice/gosip/collutil"
	. "github.com/ghettovoice/gosip/repository"
	. "github.com/ghettovoice/gosip/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
	InitTestSuite(t, "Repository Suite")
}

type behaviorInput struct {
	repo     Repository
	items    map[string]interface{}
	key      string
	expected interface{}
}

func BeEmptyBehavior(input *behaviorInput) {
	It("should be empty", func() {
		Expect(input.repo.Len()).To(Equal(0))
		Expect(input.repo.Keys()).To(HaveLen(0))
	})
}

func HaveLengthBehavior(input *behaviorInput) {
	It(fmt.Sprintf("should have length = %d", len(input.items)), func() {
		Expect(input.repo.Len()).To(Equal(len(input.items)))
		Expect(input.repo.Keys()).To(HaveLen(len(input.items)))
	})
}

func HaveItemsBehavior(input *behaviorInput) {
	It(fmt.Sprintf("should have items: %+v", input.items), func() {
		Expect(input.repo.Items()).To(ConsistOf(collutil.Vals(input.items)))
		Expect(input.repo.All()).To(ConsistOf(collutil.Vals(input.items)))
	})
	It(fmt.Sprintf("should have keys: %v", collutil.Keys(input.items)), func() {
		Expect(input.repo.Keys()).To(ConsistOf(collutil.Keys(input.items)))
	})
}

func ItemFoundBehaviour(input *behaviorInput) {
	It(fmt.Sprintf("should find item by key %s", input.key), func() {
		item, ok := input.repo.Get(input.key)
		Expect(ok).To(BeTrue())
		Expect(item).To(Equal(input.expected))
	})
}

func ItemNotFoundBehaviour(input *behaviorInput) {
	It(fmt.Sprintf("should not find item by key %s", input.key), func() {
		item, ok := input.repo.Get(input.key)
		Expect(ok).To(BeFalse())
		Expect(item).To(BeNil())
	})
}

package repository

import (
	"fmt"

	"github.com/orcaman/concurrent-map"
)

type memory struct {
	items cmap.ConcurrentMap
}

// NewMemory creates new thread-safe in-memory repository.
func NewMemory() Repository {
	repo := &memory{
		items: cmap.New(),
	}
	return repo
}

func (repo *memory) String() string {
	if repo == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%T (len: %d)", repo, repo.Len())
}

func (repo *memory) Len() int {
	return len(repo.Keys())
}

func (repo *memory) Keys() []string {
	return repo.items.Keys()
}

func (repo *memory) Has(key string) bool {
	return repo.items.Has(key)
}

func (repo *memory) Put(key string, item interface{}) {
	repo.items.Set(key, item)
}

func (repo *memory) Get(key string) (interface{}, bool) {
	return repo.items.Get(key)
}

func (repo *memory) All() map[string]interface{} {
	return repo.items.Items()
}

func (repo *memory) Drop(key string) {
	repo.items.Remove(key)
}

func (repo *memory) Clear() {
	for _, key := range repo.Keys() {
		repo.Drop(key)
	}
}

func (repo *memory) Pop(key string) (interface{}, bool) {
	return repo.items.Pop(key)
}

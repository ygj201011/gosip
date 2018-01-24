package repository

import (
	"github.com/orcaman/concurrent-map"
)

type concurrent struct {
	items cmap.ConcurrentMap
}

// NewConcurrent creates new thread-safe in-memory repository safe to use in multiple goroutines.
func NewConcurrent() Repository {
	repo := &concurrent{
		items: cmap.New(),
	}
	return repo
}

func (repo *concurrent) String() string {
	return Stringify(repo)
}

func (repo *concurrent) Len() int {
	return repo.items.Count()
}

func (repo *concurrent) Keys() []string {
	return repo.items.Keys()
}

func (repo *concurrent) Has(key string) bool {
	return repo.items.Has(key)
}

func (repo *concurrent) Put(key string, item interface{}) {
	repo.items.Set(key, item)
}

func (repo *concurrent) Get(key string) (interface{}, bool) {
	return repo.items.Get(key)
}

func (repo *concurrent) Items() map[string]interface{} {
	return repo.items.Items()
}

func (repo *concurrent) All() []interface{} {
	items := make([]interface{}, 0)
	for _, item := range repo.Items() {
		items = append(items, item)
	}
	return items
}

func (repo *concurrent) Drop(key string) {
	repo.items.Remove(key)
}

func (repo *concurrent) Clear() {
	for _, key := range repo.Keys() {
		repo.Drop(key)
	}
}

func (repo *concurrent) Pop(key string) (interface{}, bool) {
	return repo.items.Pop(key)
}

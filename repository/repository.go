// Package repository provides common types and methods to work with different repositories.
package repository

import "fmt"

// Repository introduces generic repository.
type Repository interface {
	String() string
	Len() int
	Keys() []string
	Has(key string) bool
	Put(key string, item interface{})
	Get(key string) (interface{}, bool)
	Items() map[string]interface{}
	All() []interface{}
	Drop(key string)
	Clear()
	Pop(key string) (interface{}, bool)
}

type repository struct {
	items map[string]interface{}
	keys  []string
}

// New creates simple in-memory map based repository prepared to use in a single goroutine.
func New() Repository {
	return &repository{
		items: make(map[string]interface{}),
		keys:  make([]string, 0),
	}
}

func (repo *repository) String() string {
	return Stringify(repo)
}

func (repo *repository) Len() int {
	return len(repo.Keys())
}

func (repo *repository) Keys() []string {
	return append([]string{}, repo.keys...)
}

func (repo *repository) Has(key string) bool {
	_, ok := repo.items[key]
	return ok
}

func (repo *repository) Put(key string, item interface{}) {
	repo.items[key] = item
	repo.keys = append(repo.keys, key)
}

func (repo *repository) Get(key string) (interface{}, bool) {
	item, ok := repo.items[key]
	return item, ok
}

func (repo *repository) Items() map[string]interface{} {
	items := make(map[string]interface{})
	for key, item := range repo.items {
		items[key] = item
	}
	return items
}

func (repo *repository) All() []interface{} {
	items := make([]interface{}, 0)
	for _, item := range repo.items {
		items = append(items, item)
	}
	return items
}

func (repo *repository) Drop(key string) {
	delete(repo.items, key)
	for i, k := range repo.keys {
		if k == key {
			repo.keys = append(repo.keys[:i], repo.keys[i+1:]...)
			break
		}
	}
}

func (repo *repository) Clear() {
	repo.items = make(map[string]interface{})
	repo.keys = make([]string, 0)
}

func (repo *repository) Pop(key string) (interface{}, bool) {
	if item, ok := repo.Get(key); ok {
		repo.Drop(key)
		return item, ok
	}
	return nil, false
}

func Stringify(repo Repository) string {
	if repo == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%T (len: %d)", repo, repo.Len())
}

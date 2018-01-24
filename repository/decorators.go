package repository

import (
	"time"

	"github.com/ghettovoice/gosip/timing"
)

// ItemWithTTL structs can be used to add items to repository with related TTLs
type ItemWithTTL struct {
	Item interface{}
	TTL  time.Duration
}

type itemWithTimer struct {
	item  interface{}
	timer timing.Timer
}

type withTTL struct {
	Repository
	ttl time.Duration
}

// NewWithTTL decorates any Repository with TTL feature.
// 	repo - decorated repository
// 	ttl - default TTL for all items
func NewWithTTL(repo Repository, ttl time.Duration) Repository {
	return &withTTL{repo, ttl}
}

func (repo *withTTL) String() string {
	return ToString(repo)
}

func (repo *withTTL) wrapItem(key string, item interface{}) *itemWithTimer {
	var timer timing.Timer

	if item, ok := item.(*ItemWithTTL); ok {
		timer = timing.AfterFunc(item.TTL, func() {
			repo.Drop(key)
		})
	} else {
		timer = timing.AfterFunc(repo.ttl, func() {
			repo.Drop(key)
		})
	}

	return &itemWithTimer{item, timer}
}

func (repo *withTTL) unwrapItem(item interface{}) interface{} {
	if item, ok := item.(*itemWithTimer); ok {
		return item.item
	} else {
		return item
	}
}

func (repo *withTTL) Put(key string, item interface{}) {
	// stop timer if key already exists in repository
	if foundItem, ok := repo.Repository.Get(key); ok {
		if foundItem, ok := foundItem.(*itemWithTimer); ok {
			foundItem.timer.Stop()
		}
	}

	repo.Repository.Put(key, repo.wrapItem(key, item))
}

func (repo *withTTL) Get(key string) (interface{}, bool) {
	if ttlItem, ok := repo.Repository.Get(key); ok {
		return repo.unwrapItem(ttlItem), true
	}
	return nil, false
}

func (repo *withTTL) Items() map[string]interface{} {
	items := make(map[string]interface{})
	for key, ttlItem := range repo.Repository.Items() {
		items[key] = repo.unwrapItem(ttlItem)
	}
	return items
}

func (repo *withTTL) All() []interface{} {
	items := make([]interface{}, 0, repo.Len())
	for _, item := range repo.Items() {
		items = append(items, item)
	}
	return items
}

func (repo *withTTL) Pop(key string) (interface{}, bool) {
	if item, ok := repo.Get(key); ok {
		repo.Drop(key)
		return item, ok
	}
	return nil, false
}

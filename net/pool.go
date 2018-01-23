package net

import (
	"fmt"
	"sync"

	"github.com/ghettovoice/gosip/core"
	"github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/repository"
)

// Pool is a generic network pool.
type Pool interface {
	log.LocalLogger
	core.Deferred
	repository.Repository
}

// pool provides common pool methods and variables.
type pool struct {
	log.LocalLogger
	hwg    *sync.WaitGroup
	done   chan struct{}
	errs   chan<- error
	cancel <-chan struct{}
}

func (pool *pool) String() string {
	return fmt.Sprintf("%T %p", pool, pool)
}

func (pool *pool) Done() <-chan struct{} {
	return pool.done
}

func (pool *pool) SetLog(logger log.Logger) {
	pool.LocalLogger.SetLog(logger.WithFields(map[string]interface{}{
		"pool": fmt.Sprintf("%v", pool),
	}))
}

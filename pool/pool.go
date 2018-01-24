package pool

import (
	"fmt"

	"github.com/ghettovoice/gosip/repository"
)

// Pool is a generic pool of entities.
type Pool interface {
	repository.Repository
}

type pool struct {
	repository.Repository
}

func New() Pool {
	return &pool{
		Repository: repository.NewMemory(),
	}
}

func (pool *pool) String() string {
	return fmt.Sprintf("%T %p", pool, pool)
}

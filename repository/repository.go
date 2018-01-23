// Package repository provides common types and methods to work with different repositories.
package repository

// Error is a generic repository error.
type Error interface {
	error
	// Operation indicates an operation error.
	Operation() bool
}

// OpError is an pool operation error.
type OpError struct {
	Err  error
	Op   string
	Repo Repository
}

func (*OpError) Operation() bool { return true }

func (err *OpError) Error() string {
	if err == nil {
		return "<nil>"
	}

	s := "OpError " + err.Op
	if err.Repo != nil {
		s += " (" + err.Repo.String() + ")"
	}
	s += ": " + err.Err.Error()

	return s
}

// Repository introduces generic repository.
type Repository interface {
	String() string
	Len() int
	Keys() []string
	Has(key string) bool
	Put(key string, item interface{})
	Get(key string) (interface{}, bool)
	All() map[string]interface{}
	Drop(key string)
	Clear()
	Pop(key string) (interface{}, bool)
}

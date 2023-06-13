package persistence

type Database[T any] interface {
	Add(entity T) error
	List() ([]T, error)
	Close() error
}

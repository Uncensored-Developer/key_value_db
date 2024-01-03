package storage

import "fmt"

type KeyNotFoundError struct {
	key string
}

func (k *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key %q not found in storage", k.key)
}

// Storage Interface representing the underlying storage of the database
type Storage interface {
	Set(dbIndex int, key string, value any) error
	Get(dbIndex int, key string) (any, error)
	Delete(dbIndex int, key string) error
	FetchAll(dbIndex int) <-chan [2]any
	Select(dbIndex string) (int, error)
}

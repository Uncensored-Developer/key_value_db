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
	Set(key string, value any) error
	Get(key string) (any, error)
	Delete(key string) error
}

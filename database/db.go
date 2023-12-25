package database

import "fmt"

type KeyNotFoundError struct {
	key string
}

func (k *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key '%s' not found in database", k.key)
}

type Database interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	Delete(key string) error
}

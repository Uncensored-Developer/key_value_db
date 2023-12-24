package database

type Database interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	Delete(key string) error
}

package storage

// Underlying in-memory hashmap storage
type inMemoryStorage struct {
	db map[string]any
}

func (i inMemoryStorage) Set(key string, value any) error {
	i.db[key] = value
	return nil
}

func (i inMemoryStorage) Get(key string) (any, error) {
	if val, ok := i.db[key]; ok {
		return val, nil
	} else {
		return nil, &KeyNotFoundError{key: key}
	}
}

func (i inMemoryStorage) Delete(key string) error {
	if _, ok := i.db[key]; ok {
		delete(i.db, key)
		return nil
	} else {
		return &KeyNotFoundError{key: key}
	}
}

func NewInMemoryStorage() Storage {
	return &inMemoryStorage{
		db: make(map[string]any),
	}
}

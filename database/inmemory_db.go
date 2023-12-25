package database

type inMemoryDB struct {
	db map[string]any
}

func (i inMemoryDB) Set(key string, value any) error {
	i.db[key] = value
	return nil
}

func (i inMemoryDB) Get(key string) (any, error) {
	if val, ok := i.db[key]; ok {
		return val, nil
	} else {
		return nil, &KeyNotFoundError{key: key}
	}
}

func (i inMemoryDB) Delete(key string) error {
	if _, ok := i.db[key]; ok {
		delete(i.db, key)
		return nil
	} else {
		return &KeyNotFoundError{key: key}
	}
}

func NewInMemoryDb() Database {
	return &inMemoryDB{
		db: make(map[string]any),
	}
}

package storage

import (
	"errors"
	"fmt"
	"strconv"
)

// Underlying in-memory hashmap storage
type inMemoryStorage struct {
	dbCount int // Number of available databases
	db      map[int]map[string]any
}

func (i inMemoryStorage) Set(dbIndex int, key string, value any) error {
	i.db[dbIndex][key] = value
	return nil
}

func (i inMemoryStorage) Get(dbIndex int, key string) (any, error) {
	if val, ok := i.db[dbIndex][key]; ok {
		return val, nil
	} else {
		return nil, &KeyNotFoundError{key: key}
	}
}

func (i inMemoryStorage) Delete(dbIndex int, key string) error {
	if _, ok := i.db[dbIndex][key]; ok {
		delete(i.db[dbIndex], key)
		return nil
	} else {
		return &KeyNotFoundError{key: key}
	}
}

func (i inMemoryStorage) FetchAll(dbIndex int) <-chan [2]any {
	outputChan := make(chan [2]any)
	go func() {
		for k, v := range i.db[dbIndex] {
			outputChan <- [2]any{k, v}
		}
		close(outputChan)
	}()
	return outputChan
}

func (i inMemoryStorage) Select(dbIndex string) (int, error) {
	dbIndexInt, err := strconv.Atoi(dbIndex)
	if err != nil {
		return 0, fmt.Errorf("(error) ERR value is not an integer or out of range")
	}
	if dbIndexInt < 0 || dbIndexInt > i.dbCount-1 {
		return 0, errors.New("(error) ERR DB index is out of range")
	}
	return dbIndexInt, nil
}

func NewInMemoryStorage(dbCount int) Storage {
	if dbCount == 0 {
		dbCount = 16
	}
	db := make(map[int]map[string]any)

	for i := 0; i < dbCount; i++ {
		db[i] = make(map[string]any)
	}
	return &inMemoryStorage{
		dbCount: dbCount,
		db:      db,
	}
}

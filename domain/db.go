package domain

import (
	"fmt"
	"kvdb/storage"
	"strconv"
)

type KeyValueDB struct {
	storage storage.Storage
}

func NewKeyValueDB(storage storage.Storage) KeyValueDB {
	return KeyValueDB{storage: storage}
}

type DBResult struct {
	Value    any
	Type     string
	Response string
}

func (k *KeyValueDB) Execute(cmd Command) (DBResult, error) {
	_, err := cmd.Validate()
	if err != nil {
		return DBResult{Value: err.Error(), Response: ""}, err
	}
	switch cmd.Keyword {
	case SET:
		err := k.storage.Set(cmd.Key, cmd.Value)
		if err != nil {
			return DBResult{Value: err.Error()}, err
		}
		return DBResult{Value: "", Response: "OK"}, nil
	case GET:
		result, err := k.storage.Get(cmd.Key)
		if err != nil {
			return DBResult{Value: err.Error(), Response: "(nil)"}, err
		}
		return DBResult{Value: result, Response: ""}, nil
	case DEL:
		err := k.storage.Delete(cmd.Key)
		if err != nil {
			return DBResult{Value: err.Error(), Type: "integer", Response: "0"}, err
		}
		return DBResult{Value: "", Type: "integer", Response: "1"}, nil
	case INCR, INCRBY:
		result, err := k.storage.Get(cmd.Key)
		if err != nil {
			return DBResult{Value: err.Error(), Response: "(nil)"}, err
		}
		intValue, err := convertToInt(result)
		if err != nil {
			return DBResult{Value: err.Error(), Response: ""}, err
		}
		change := 1
		if cmd.Keyword == INCRBY {
			intSetValue, err := convertToInt(cmd.Value)
			if err != nil {
				return DBResult{Value: err.Error(), Response: ""}, err
			}
			change = intSetValue
		}
		newValue := intValue + change
		err = k.storage.Set(cmd.Key, newValue)
		if err != nil {
			return DBResult{Value: err.Error()}, err
		}
		return DBResult{Value: newValue, Type: "integer", Response: ""}, nil
	}
	return DBResult{}, nil
}

func convertToInt(value any) (int, error) {
	intValue, err := strconv.Atoi(value.(string))
	if err != nil {
		return 0, fmt.Errorf("(error) ERR value is not an integer")
	}
	return intValue, nil
}

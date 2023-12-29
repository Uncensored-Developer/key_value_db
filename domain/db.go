package domain

import "kvdb/storage"

type KeyValueDB struct {
	storage storage.Storage
}

func NewKeyValueDB(storage storage.Storage) KeyValueDB {
	return KeyValueDB{storage: storage}
}

type DBResult struct {
	Value    any
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
			return DBResult{Value: err.Error(), Response: "0"}, err
		}
		return DBResult{Value: "", Response: "1"}, nil
	}
	return DBResult{}, nil
}

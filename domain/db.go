package domain

import "kvdb/storage"

type KeyValueDB struct {
	storage storage.Storage
}

func NewKeyValueDB(storage storage.Storage) KeyValueDB {
	return KeyValueDB{storage: storage}
}

func (k *KeyValueDB) Execute(cmd Command) (any, error) {
	_, err := cmd.Validate()
	if err != nil {
		return "", err
	}
	switch cmd.Keyword {
	case SET:
		err := k.storage.Set(cmd.Key, cmd.Value)
		if err != nil {
			return "", err
		}
		return "OK", nil
	case GET:
		result, err := k.storage.Get(cmd.Key)
		if err != nil {
			return "(nil)", err
		}
		return result, nil
	case DEL:
		err := k.storage.Delete(cmd.Key)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return "", nil
}

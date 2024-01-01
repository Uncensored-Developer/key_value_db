package domain

import (
	"fmt"
	"kvdb/storage"
	"reflect"
	"strconv"
	"strings"
)

type KeyValueDB struct {
	storage            storage.Storage
	cmdQueue           []Command
	multiCommandActive bool
}

func NewKeyValueDB(storage storage.Storage) KeyValueDB {
	return KeyValueDB{storage: storage}
}

type DBResult struct {
	Value    any
	Type     string
	Response string
	Err      error
}

type MultiBlockError struct {
	cmd string
}

func (m *MultiBlockError) Error() string {
	return fmt.Sprintf("(error) ERR %s without MULTI", m.cmd)
}

func (d DBResult) SimpleMsg() any {
	value := d.Value
	err := d.Err

	if d.Response == "" {
		if isString(d.Value) && err == nil {
			value = fmt.Sprintf("%q", d.Value)
		}
	} else {
		value = d.Response
	}
	if d.Type != "" {
		value = fmt.Sprintf("(%s) %v", d.Type, value)
	}
	return value
}

func (d DBResult) String() string {
	return fmt.Sprintf("{Value: %v, Type: %q, Response: %q, Err: %v}", d.Value, d.Type, d.Response, d.Err)
}

func (k *KeyValueDB) Execute(cmd Command) any {
	_, err := cmd.Validate()
	if err != nil {
		return DBResult{Value: err.Error(), Response: "", Err: err}
	}

	if k.multiCommandActive && !cmd.isExitMultiBlockCmd() {
		k.cmdQueue = append(k.cmdQueue, cmd)
		return DBResult{Value: "", Response: "QUEUED"}
	}

	switch cmd.Keyword {
	case SET:
		err := k.storage.Set(cmd.Key, cmd.Value)
		if err != nil {
			return DBResult{Value: err.Error(), Err: err}
		}
		return DBResult{Value: "", Response: "OK"}
	case GET:
		result, err := k.storage.Get(cmd.Key)
		if err != nil {
			return DBResult{Value: err.Error(), Response: "(nil)", Err: err}
		}
		return DBResult{Value: result, Response: ""}
	case DEL:
		err := k.storage.Delete(cmd.Key)
		if err != nil {
			return DBResult{Value: err.Error(), Type: "integer", Response: "0", Err: err}
		}
		return DBResult{Value: "", Type: "integer", Response: "1"}
	case INCR, INCRBY:
		result, err := k.storage.Get(cmd.Key)
		if err != nil {
			return DBResult{Value: err.Error(), Response: "(nil)", Err: err}
		}
		intValue, err := convertToInt(result)
		if err != nil {
			return DBResult{Value: err.Error(), Response: "", Err: err}
		}
		change := 1
		if cmd.Keyword == INCRBY {
			intSetValue, err := convertToInt(cmd.Value)
			if err != nil {
				return DBResult{Value: err.Error(), Response: "", Err: err}
			}
			change = intSetValue
		}
		newValue := intValue + change
		err = k.storage.Set(cmd.Key, newValue)
		if err != nil {
			return DBResult{Value: err.Error(), Err: err}
		}

		return DBResult{Value: newValue, Type: "integer", Response: ""}
	case MULTI:
		k.multiCommandActive = true
		return DBResult{Value: "", Response: "OK"}
	case DISCARD:
		if !k.multiCommandActive {
			err = &MultiBlockError{cmd: DISCARD}
			return DBResult{Value: err.Error(), Response: "", Err: err}
		}
		k.multiCommandActive = false
		k.cmdQueue = nil
		return DBResult{Value: "", Response: "OK"}
	case EXEC:
		if !k.multiCommandActive {
			err = &MultiBlockError{cmd: EXEC}
			return DBResult{Value: err.Error(), Response: "", Err: err}
		}
		k.multiCommandActive = false
		return k.executeQueuedCmds()
	case COMPACT:
		var results []DBResult
		for kv := range k.storage.FetchAll() {
			cmdKey := kv[0].(string)
			value := kv[1]

			if len(strings.Fields(cmdKey)) > 1 {
				cmdKey = fmt.Sprintf("%q", cmdKey)
			}

			switch v := value.(type) {
			case int:
				value = v
			case string:
				value = fmt.Sprintf("%q", value)
			}

			dbRes := DBResult{Response: fmt.Sprintf("SET %s %v", cmdKey, value)}
			results = append(results, dbRes)
		}
		return results
	}
	return DBResult{}
}

// executeQueuedCmds executes the queued commands in the KeyValueDB.
//
// It iterates over the cmdQueue and executes each command using the Execute method of KeyValueDB. The results of each execution are stored in the results slice. After executing all the commands, the cmdQueue is set to nil. The function then returns the results slice.
// Returns []ExecuteReturnValue.
func (k *KeyValueDB) executeQueuedCmds() []DBResult {
	var results []DBResult
	for _, cmd := range k.cmdQueue {
		dbRes := k.Execute(cmd)
		results = append(results, dbRes.(DBResult))
	}
	k.cmdQueue = nil
	return results
}

func convertToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case string:
		intValue, err := strconv.Atoi(value.(string))
		if err != nil {
			return 0, fmt.Errorf("(error) ERR value is not an integer")
		}
		return intValue, nil
	default:
		return 0, fmt.Errorf("(error) ERR value is neither string nor int")
	}
}

func isString(obj any) bool {
	return reflect.TypeOf(obj).Kind() == reflect.String
}

package domain

import (
	"errors"
	"kvdb/storage"
	"reflect"
	"testing"
)

func TestKeyValueDB_Execute(t *testing.T) {
	testCases := []struct {
		name        string
		cmds        []Command
		wantResults []any
		wantErrMsgs []string
	}{
		{
			name:        "Valid SET command",
			cmds:        []Command{NewCommand("SET", "key", "value")},
			wantResults: []any{"OK"},
			wantErrMsgs: []string{""},
		},
		{
			name:        "Invalid SET command",
			cmds:        []Command{NewCommand("SET", "key")},
			wantResults: []any{"(error) ERR SET command expected 2 arguments but 1 was given (i.e no value)"},
			wantErrMsgs: []string{"(error) ERR SET command expected 2 arguments but 1 was given (i.e no value)"},
		},
		{
			name:        "Get non-existing key",
			cmds:        []Command{NewCommand("GET", "non-existing_key")},
			wantResults: []any{"(nil)"},
			wantErrMsgs: []string{"Key \"non-existing_key\" not found in storage"},
		},
		{
			name: "Get existing key",
			cmds: []Command{
				NewCommand("SET", "key", "value"),
				NewCommand("GET", "key"),
			},
			wantResults: []any{"OK", "value"},
			wantErrMsgs: []string{"", ""},
		},
		{
			name:        "Delete non-existing key",
			cmds:        []Command{NewCommand("DEL", "non-existing_key")},
			wantResults: []any{"0"},
			wantErrMsgs: []string{"Key \"non-existing_key\" not found in storage"},
		},
		{
			name: "Delete existing key",
			cmds: []Command{
				NewCommand("SET", "key", "value"),
				NewCommand("DEL", "key"),
			},
			wantResults: []any{"OK", "1"},
			wantErrMsgs: []string{"", ""},
		},
		{
			name:        "Increase non-existing key",
			cmds:        []Command{NewCommand("INCR", "non-existing_key")},
			wantResults: []any{"(nil)"},
			wantErrMsgs: []string{"Key \"non-existing_key\" not found in storage"},
		},
		{
			name: "Increase non-integer value",
			cmds: []Command{
				NewCommand("SET", "key", "abc"),
				NewCommand("INCR", "key"),
			},
			wantResults: []any{"OK", "(error) ERR value is not an integer"},
			wantErrMsgs: []string{"", "(error) ERR value is not an integer"},
		},
		{
			name: "Increase valid-integer value",
			cmds: []Command{
				NewCommand("SET", "key", "10"),
				NewCommand("INCR", "key"),
			},
			wantResults: []any{"OK", 11},
			wantErrMsgs: []string{"", ""},
		},
		{
			name:        "IncreaseBy non-existing key",
			cmds:        []Command{NewCommand("INCRBY", "non-existing_key", "10")},
			wantResults: []any{"(nil)"},
			wantErrMsgs: []string{"Key \"non-existing_key\" not found in storage"},
		},
		{
			name: "IncreaseBy non-integer value",
			cmds: []Command{
				NewCommand("SET", "key", "10"),
				NewCommand("INCRBY", "key", "abc"),
			},
			wantResults: []any{"OK", "(error) ERR value is not an integer"},
			wantErrMsgs: []string{"", "(error) ERR value is not an integer"},
		},
		{
			name: "IncreaseBy valid-integer value",
			cmds: []Command{
				NewCommand("SET", "key", "10"),
				NewCommand("INCRBY", "key", "5"),
			},
			wantResults: []any{"OK", 15},
			wantErrMsgs: []string{"", ""},
		},
		{
			name: "Multi",
			cmds: []Command{
				NewCommand("MULTI"),
				NewCommand("SET", "key", "5"),
				NewCommand("INCR", "key"),
				NewCommand("INCRBY", "key", "5"),
			},
			wantResults: []any{"OK", "QUEUED", "QUEUED", "QUEUED"},
			wantErrMsgs: []string{"", "", "", ""},
		},
		{
			name:        "Discard - without MULTI block",
			cmds:        []Command{NewCommand("DISCARD")},
			wantResults: []any{"(error) ERR DISCARD without MULTI"},
			wantErrMsgs: []string{"(error) ERR DISCARD without MULTI"},
		},
		{
			name: "Discard - with MULTI block",
			cmds: []Command{
				NewCommand("MULTI"),
				NewCommand("SET", "key", "5"),
				NewCommand("INCR", "key"),
				NewCommand("INCRBY", "key", "5"),
				NewCommand("DISCARD"),
			},
			wantResults: []any{"OK", "QUEUED", "QUEUED", "QUEUED", "OK"},
			wantErrMsgs: []string{"", "", "", "", ""},
		},
		{
			name:        "Exec - without MULTI block",
			cmds:        []Command{NewCommand("EXEC")},
			wantResults: []any{"(error) ERR EXEC without MULTI"},
			wantErrMsgs: []string{"(error) ERR EXEC without MULTI"},
		},
	}

	for _, tc := range testCases {
		inMemoryStorage := storage.NewInMemoryStorage()
		db := NewKeyValueDB(inMemoryStorage)

		t.Run(tc.name, func(t *testing.T) {

			for i, cmd := range tc.cmds {
				got := db.Execute(cmd).(DBResult)
				gotErr := got.Err

				if gotErr == nil {
					// Placeholder for testing nil errors
					gotErr = errors.New("")
				}

				if gotErr.Error() != tc.wantErrMsgs[i] {
					t.Fatalf("KeyValueDB.Execute(%v) = %v, want Error %v", tc.cmds[i], gotErr, tc.wantErrMsgs[i])
				}

				if got.Response == "" {
					if got.Value != tc.wantResults[i] {
						t.Errorf("KeyValueDB.Execute(%v) = %q, want %q", tc.cmds[i], got.Value, tc.wantResults[i])
					}
				} else {
					if got.Response != tc.wantResults[i] {
						t.Errorf("KeyValueDB.Execute(%v) = %q, want %q", tc.cmds[i], got.Response, tc.wantResults[i])
					}
				}

			}

		})
	}
}

func TestKeyValueDB_Execute_ExecCommand(t *testing.T) {
	inMemoryStorage := storage.NewInMemoryStorage()
	db := NewKeyValueDB(inMemoryStorage)

	want := []DBResult{
		{Value: "", Response: "OK"},
		{Value: 6, Type: "integer"},
		{Value: 11, Type: "integer"},
	}

	var cmds []Command = []Command{
		NewCommand("MULTI"),
		NewCommand("SET", "key", "5"),
		NewCommand("INCR", "key"),
		NewCommand("INCRBY", "key", "5"),
	}

	for _, cmd := range cmds {
		got := db.Execute(cmd).(DBResult)
		gotErr := got.Err
		if gotErr != nil {
			t.Fatalf("Unexpected error: %v", gotErr)
		}
	}

	got := db.Execute(NewCommand("EXEC")).([]DBResult)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("EXEC command got %v, want %v", got, want)
	}

}

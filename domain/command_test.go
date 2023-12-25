package domain

import (
	"reflect"
	"testing"
)

func TestNewCommand(t *testing.T) {

	testCases := []struct {
		name    string
		keyword string
		args    []any
		want    Command
	}{
		{
			name:    "Empty keyword",
			keyword: "",
			args:    []any{},
			want:    Command{},
		},
		{
			name:    "Keyword with one argument",
			keyword: GET,
			args:    []any{"key"},
			want:    Command{Keyword: "GET", key: "key"},
		},
		{
			name:    "Keyword with two argument",
			keyword: SET,
			args:    []any{"key", "value"},
			want:    Command{Keyword: "SET", key: "key", Value: "value"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewCommand(tc.keyword, tc.args...)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewCommand() = %v, want %v", got, tc.want)
			}
		})
	}

}

func TestCommand_Validate(t *testing.T) {

	testCases := []struct {
		name          string
		command       Command
		wantValidated bool
		wantError     error
	}{
		{
			name:          "Empty keyword command",
			command:       Command{},
			wantValidated: false,
			wantError:     &CommandError{msg: "unknown command "},
		},
		{
			name:          "Invalid keyword command",
			command:       Command{Keyword: "PUT"},
			wantValidated: false,
			wantError:     &CommandError{msg: "unknown command PUT"},
		},
		{
			name:          "SET command - no key and value",
			command:       Command{Keyword: "SET"},
			wantValidated: false,
			wantError:     &CommandError{msg: "SET command expected 2 arguments but none was given (i.e no key & value)"},
		},
		{
			name:          "SET command - key and no value",
			command:       Command{Keyword: "SET", key: "key_1"},
			wantValidated: false,
			wantError:     &CommandError{msg: "SET command expected 2 arguments but 1 was given (i.e no value)"},
		},
		{
			name:          "SET command - valid key and value",
			command:       Command{Keyword: "SET", key: "key_1", Value: "value_1"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "GET command - no key",
			command:       Command{Keyword: "GET"},
			wantValidated: false,
			wantError:     &CommandError{msg: "GET command expected 1 argument but none was given (i.e no key)"},
		},
		{
			name:          "GET command - key and value",
			command:       Command{Keyword: "GET", key: "key_1", Value: "value_1"},
			wantValidated: false,
			wantError:     &CommandError{msg: "GET command expected 1 argument but 2 was given"},
		},
		{
			name:          "GET command - valid key",
			command:       Command{Keyword: "GET", key: "key_1"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "DEL command - valid key",
			command:       Command{Keyword: "DEL", key: "key_1"},
			wantValidated: true,
			wantError:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.command.Validate()
			if got != tc.wantValidated {
				t.Errorf("Command.Validate() = %v, want %v", got, tc.wantValidated)
			}

			if !reflect.DeepEqual(err, tc.wantError) {
				t.Errorf("Command.Validate() = %v, want Error %v", err, tc.wantError)
			}
		})
	}

}

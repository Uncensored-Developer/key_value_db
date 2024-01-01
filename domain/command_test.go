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
			args:    []any{"Key"},
			want:    Command{Keyword: "GET", Key: "Key"},
		},
		{
			name:    "Keyword with two argument",
			keyword: SET,
			args:    []any{"Key", "value"},
			want:    Command{Keyword: "SET", Key: "Key", Value: "value"},
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
			name:          "SET command - no Key and value",
			command:       Command{Keyword: "SET"},
			wantValidated: false,
			wantError:     &CommandError{msg: "SET command expected 2 arguments but none was given (i.e no Key & value)"},
		},
		{
			name:          "SET command - Key and no value",
			command:       Command{Keyword: "SET", Key: "key_1"},
			wantValidated: false,
			wantError:     &CommandError{msg: "SET command expected 2 arguments but 1 was given (i.e no value)"},
		},
		{
			name:          "SET command - valid Key and value",
			command:       Command{Keyword: "SET", Key: "key_1", Value: "value_1"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "GET command - no Key",
			command:       Command{Keyword: "GET"},
			wantValidated: false,
			wantError:     &CommandError{msg: "GET command expected 1 argument but none was given (i.e no Key)"},
		},
		{
			name:          "GET command - Key and value",
			command:       Command{Keyword: "GET", Key: "key_1", Value: "value_1"},
			wantValidated: false,
			wantError:     &CommandError{msg: "GET command expected 1 argument but 2 was given"},
		},
		{
			name:          "GET command - valid Key",
			command:       Command{Keyword: "GET", Key: "key_1"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "DEL command - valid Key",
			command:       Command{Keyword: "DEL", Key: "key_1"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "INCR command - valid Key",
			command:       Command{Keyword: "INCR", Key: "key_1"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "INCRBY command - valid Key",
			command:       Command{Keyword: "INCRBY", Key: "key_1", Value: 10},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "MULTI command - Key and value",
			command:       Command{Keyword: "MULTI", Key: "key_1", Value: "value_1"},
			wantValidated: false,
			wantError:     &CommandError{msg: "MULTI command expected no argument but was given"},
		},
		{
			name:          "MULTI command - valid",
			command:       Command{Keyword: "MULTI"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "DISCARD command - valid",
			command:       Command{Keyword: "DISCARD"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "EXEC command - valid",
			command:       Command{Keyword: "EXEC"},
			wantValidated: true,
			wantError:     nil,
		},
		{
			name:          "COMPACT command - valid",
			command:       Command{Keyword: "COMPACT"},
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

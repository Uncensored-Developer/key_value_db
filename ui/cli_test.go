package ui

import (
	"errors"
	"kvdb/domain"
	"testing"
)

func Test_getCommand(t *testing.T) {

	testCases := []struct {
		name       string
		input      string
		want       domain.Command
		wantErrMsg string
	}{
		{
			name:       "SET command - one word value",
			input:      "SET key value",
			want:       domain.Command{Keyword: "SET", Key: "key", Value: "value"},
			wantErrMsg: "",
		},
		{
			name:       "SET command - argument with extra spaces",
			input:      " SET  key  value ",
			want:       domain.Command{Keyword: "SET", Key: "key", Value: "value"},
			wantErrMsg: "",
		},
		{
			name:       "GET command - one word key",
			input:      " GET key",
			want:       domain.Command{Keyword: "GET", Key: "key"},
			wantErrMsg: "",
		},
		{
			name:       "GET command - lowercase keyword",
			input:      " get key",
			want:       domain.Command{Keyword: "GET", Key: "key"},
			wantErrMsg: "",
		},
		{
			name:       "GET command - multiword key in quotes",
			input:      "GET \"multi key\"",
			want:       domain.Command{Keyword: "GET", Key: "multi key"},
			wantErrMsg: "",
		},
		{
			name:       "GET command - multiword key with extra spaces in quotes",
			input:      "GET \"multi  key\"",
			want:       domain.Command{Keyword: "GET", Key: "multi  key"},
			wantErrMsg: "",
		},
		{
			name:       "GET command - 5 word key with extra spaces in quotes",
			input:      "GET \"one two  three four  key\"",
			want:       domain.Command{Keyword: "GET", Key: "one two  three four  key"},
			wantErrMsg: "",
		},
		{
			name:       "SET command - multiword key and value in quotes",
			input:      "SET \"multi word key\" \"multi word value\"",
			want:       domain.Command{Keyword: "SET", Key: "multi word key", Value: "multi word value"},
			wantErrMsg: "",
		},
		{
			name:       "SET command - no closing quote",
			input:      "SET \"multi word key\" \"multi word value",
			want:       domain.Command{},
			wantErrMsg: "(error) ERR Syntax error: arguments has no closing quote",
		},
		{
			name:       "SET command - syntax error",
			input:      "SET key value1 value2",
			want:       domain.Command{},
			wantErrMsg: "(error) ERR Syntax error",
		},
		{
			name:       "SET command - syntax error - too many arguments",
			input:      "SET \"multi word key\" \"multi word value1\" \"multi word value2\"",
			want:       domain.Command{},
			wantErrMsg: "(error) ERR Syntax error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			got, gotErr := getCommand(tc.input)

			if gotErr == nil {
				// Placeholder for testing nil errors
				gotErr = errors.New("")
			}

			if gotErr.Error() != tc.wantErrMsg {
				t.Fatalf("getCommand(%q) = %v, want Error %v", tc.input, gotErr, tc.wantErrMsg)
			}

			if got != tc.want {
				t.Errorf("getCommand(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}

}

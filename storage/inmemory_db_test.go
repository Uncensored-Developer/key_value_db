package storage

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewInMemoryStorage(t *testing.T) {
	testCases := []struct {
		name    string
		dbCount int
		want    Storage
	}{
		{
			name:    "Zero dbCount should default to 16",
			dbCount: 0,
			want:    NewInMemoryStorage(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewInMemoryStorage(tc.dbCount)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewInMemory() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInMemoryDB_SetGet(t *testing.T) {

	t.Run("Set a key-value pair", func(t *testing.T) {
		db := NewInMemoryStorage(0)
		dbIndex := 0
		key := "key_1"
		want := "value_1"
		err := db.Set(dbIndex, key, "value_1")
		if err != nil {
			t.Fatalf("Unexpected error setting key-value pair: %v", err)
		}

		got, err := db.Get(dbIndex, key)
		if err != nil {
			t.Errorf("Unexpected error getting value from '%s': %v", key, err)
		}

		if got != want {
			t.Errorf("inMemory.Set() = %v, want %v", got, want)
		}
	})

	t.Run("Update the value of exising key", func(t *testing.T) {
		db := NewInMemoryStorage(0)
		dbIndex := 0
		key := "key_1"
		want := "value_2"
		err := db.Set(dbIndex, key, "value_1")
		if err != nil {
			t.Fatalf("Unexpected error setting key-value pair: %v", err)
		}

		_ = db.Set(dbIndex, key, "value_2")

		got, _ := db.Get(dbIndex, key)

		if got != want {
			t.Errorf("inMemory.Set() = %v, want %v", got, want)
		}
	})

	t.Run("Get the value of non-exising key", func(t *testing.T) {
		db := NewInMemoryStorage(0)
		dbIndex := 0
		key := "invalid_key"
		wantErr := KeyNotFoundError{key: key}

		_, gotErr := db.Get(dbIndex, key)

		if gotErr == nil {
			t.Fatalf("inMemory.Get() expected an error, but got nil")
		}

		var notFoundErr *KeyNotFoundError
		if !errors.As(gotErr, &notFoundErr) {
			t.Errorf("Got error %s, want %s", gotErr.Error(), wantErr.Error())
		}
	})

}

func TestInMemoryDB_Delete(t *testing.T) {
	t.Run("Delete a non-existing key", func(t *testing.T) {
		db := NewInMemoryStorage(0)
		dbIndex := 0
		key := "invalid_key"
		wantErr := KeyNotFoundError{key: key}

		gotErr := db.Delete(dbIndex, key)

		if gotErr == nil {
			t.Fatalf("inMemory.Delete() expected an error, but got nil")
		}

		var notFoundErr *KeyNotFoundError
		if !errors.As(gotErr, &notFoundErr) {
			t.Errorf("Got error %s, want %s", gotErr.Error(), wantErr.Error())
		}
	})

	t.Run("Delete an existing key", func(t *testing.T) {
		db := NewInMemoryStorage(0)
		dbIndex := 0
		key := "key_3"
		value := "value_3"

		_ = db.Set(dbIndex, key, value)

		_, getErr := db.Get(dbIndex, key)
		if getErr != nil {
			t.Errorf("Unexpected error getting value from with key '%s': %v", key, getErr)
		}

		delErr := db.Delete(dbIndex, key)
		if delErr != nil {
			t.Fatalf("inMemory.Delete() = %v, want %v", delErr, nil)
		}

		_, getErr = db.Get(dbIndex, key)
		wantErr := KeyNotFoundError{key: key}
		if getErr == nil {
			t.Fatalf("inMemory.Get() expected an error, but got nil")
		}

		var notFoundErr *KeyNotFoundError
		if !errors.As(getErr, &notFoundErr) {
			t.Errorf("Got error %s, want %s", getErr.Error(), wantErr.Error())
		}
	})
}

func TestInMemoryStorage_Select(t *testing.T) {
	testCases := []struct {
		name    string
		dbCount int
		dbIndex string
		want    int
		wantErr error
	}{
		{
			name:    "Invalid dbIndex",
			dbCount: 10,
			dbIndex: "invalid",
			want:    0,
			wantErr: errors.New("(error) ERR value is not an integer or out of range"),
		},
		{
			name:    "Invalid dbIndex - out of range",
			dbCount: 10,
			dbIndex: "16",
			want:    0,
			wantErr: errors.New("(error) ERR DB index is out of range"),
		},
		{
			name:    "Valid dbIndex",
			dbCount: 10,
			dbIndex: "3",
			want:    3,
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryStorage(tc.dbCount)
			got, err := db.Select(tc.dbIndex)
			if err != nil {
				if tc.wantErr == nil {
					t.Errorf("inMemory.Select() returned expected error: %v", err)
				} else if err.Error() != tc.wantErr.Error() {
					t.Errorf("inMemory.Select() error = %v, want error %v", got, tc.wantErr)
				}
			} else if tc.wantErr != nil {
				t.Errorf("inMemory.Select() error = <nil>, want error %v", tc.wantErr)
			}

			if got != tc.want {
				t.Errorf("inMemory.Select() = %v, want %v", got, tc.want)
			}
		})
	}
}

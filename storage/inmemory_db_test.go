package storage

import (
	"errors"
	"testing"
)

func TestInMemoryDB_SetGet(t *testing.T) {

	t.Run("Set a key-value pair", func(t *testing.T) {
		db := NewInMemoryDb()
		key := "key_1"
		want := "value_1"
		err := db.Set(key, "value_1")
		if err != nil {
			t.Fatalf("Unexpected error setting key-value pair: %v", err)
		}

		got, err := db.Get(key)
		if err != nil {
			t.Errorf("Unexpected error getting value from '%s': %v", key, err)
		}

		if got != want {
			t.Errorf("inMemory.Set() = %v, want %v", got, want)
		}
	})

	t.Run("Update the value of exising key", func(t *testing.T) {
		db := NewInMemoryDb()
		key := "key_1"
		want := "value_2"
		err := db.Set(key, "value_1")
		if err != nil {
			t.Fatalf("Unexpected error setting key-value pair: %v", err)
		}

		_ = db.Set(key, "value_2")

		got, _ := db.Get(key)

		if got != want {
			t.Errorf("inMemory.Set() = %v, want %v", got, want)
		}
	})

	t.Run("Get the value of non-exising key", func(t *testing.T) {
		db := NewInMemoryDb()
		key := "invalid_key"
		wantErr := KeyNotFoundError{key: key}

		_, gotErr := db.Get(key)

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
		db := NewInMemoryDb()
		key := "invalid_key"
		wantErr := KeyNotFoundError{key: key}

		gotErr := db.Delete(key)

		if gotErr == nil {
			t.Fatalf("inMemory.Delete() expected an error, but got nil")
		}

		var notFoundErr *KeyNotFoundError
		if !errors.As(gotErr, &notFoundErr) {
			t.Errorf("Got error %s, want %s", gotErr.Error(), wantErr.Error())
		}
	})

	t.Run("Delete an existing key", func(t *testing.T) {
		db := NewInMemoryDb()
		key := "key_3"
		value := "value_3"

		_ = db.Set(key, value)

		_, getErr := db.Get(key)
		if getErr != nil {
			t.Errorf("Unexpected error getting value from with key '%s': %v", key, getErr)
		}

		delErr := db.Delete(key)
		if delErr != nil {
			t.Fatalf("inMemory.Delete() = %v, want %v", delErr, nil)
		}

		_, getErr = db.Get(key)
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

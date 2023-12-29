package main

import (
	"kvdb/domain"
	"kvdb/storage"
	"kvdb/ui"
)

func main() {
	inMemoryStorage := storage.NewInMemoryStorage()
	keyValueDB := domain.NewKeyValueDB(inMemoryStorage)
	ui.RunCLI(keyValueDB)
}

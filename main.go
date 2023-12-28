package main

import (
	"kvdb/storage"
	"kvdb/ui"
)

func main() {
	keyValueDB := storage.NewInMemoryStorage()
	ui.RunCLI(keyValueDB)
}

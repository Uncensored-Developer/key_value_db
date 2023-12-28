package main

import (
	"kvdb/storage"
	"kvdb/ui"
)

func main() {
	keyValueDB := storage.NewInMemoryDb()
	ui.RunCLI(keyValueDB)
}

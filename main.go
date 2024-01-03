package main

import (
	"github.com/joho/godotenv"
	"kvdb/domain"
	"kvdb/storage"
	"kvdb/ui"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var shutDownSignal = make(chan os.Signal, 1)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("TCP_PORT")

	inMemoryStorage := storage.NewInMemoryStorage()
	keyValueDB := domain.NewKeyValueDB(inMemoryStorage)

	tcpServer := ui.NewTcpServer(port, keyValueDB)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal

	tcpServer.Stop()
}

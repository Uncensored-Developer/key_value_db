package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"kvdb/domain"
	"kvdb/storage"
	"kvdb/ui"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var shutDownSignal = make(chan os.Signal, 1)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("TCP_PORT")
	dbCount := os.Getenv("DB_COUNT")

	dbCountInt, err := getIntDbCount(dbCount)
	if err != nil {
		log.Fatalf("Error setting DB_COUNT: %v", err)
	}

	inMemoryStorage := storage.NewInMemoryStorage(dbCountInt)
	keyValueDB := domain.NewKeyValueDB(inMemoryStorage)

	tcpServer := ui.NewTcpServer(port, keyValueDB)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal

	tcpServer.Stop()
}

func getIntDbCount(dbCountStr string) (int, error) {
	dbCountInt, err := strconv.Atoi(dbCountStr)
	if err != nil {
		return 0, fmt.Errorf("error converting dbCountStr to int: %v", err)
	}
	return dbCountInt, nil
}

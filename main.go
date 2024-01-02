package main

import (
	"github.com/joho/godotenv"
	"kvdb/domain"
	"kvdb/storage"
	"kvdb/ui"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("TCP_PORT")

	inMemoryStorage := storage.NewInMemoryStorage()
	keyValueDB := domain.NewKeyValueDB(inMemoryStorage)

	listener, err := ui.StartTcpServer(port)
	if err != nil {
		log.Fatalf("Failed to startup TCP server: %v\n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v\n", err)
		}
		go ui.HandleConnection(conn, keyValueDB)
	}
}

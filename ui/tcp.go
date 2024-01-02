package ui

import (
	"bufio"
	"fmt"
	"kvdb/domain"
	"log"
	"net"
)

func HandleConnection(conn net.Conn, db domain.KeyValueDB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			PrintDbResult(writer, domain.DBResult{Value: err.Error(), Err: err})
			break
		}

		command, err := getCommand(input)
		if err != nil {
			PrintDbResult(writer, err)
			break
		}
		result := db.Execute(command)
		PrintDbResult(writer, result)
		err = writer.Flush()
		if err != nil {
			log.Printf("Error flusing buffered writer: %v\n", err)
			return
		}
	}
}

func StartTcpServer(port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Server listening on port %s\n", port)
	return listener, nil
}

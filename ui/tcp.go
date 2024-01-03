package ui

import (
	"bufio"
	"fmt"
	"kvdb/domain"
	"log"
	"net"
	"sync"
)

type TcpServer struct {
	listener net.Listener
	shutdown chan struct{}
	wg       sync.WaitGroup
	db       domain.KeyValueDB
}

func NewTcpServer(port string, db domain.KeyValueDB) *TcpServer {
	s := &TcpServer{
		shutdown: make(chan struct{}),
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to startup TCP server: %v\n", err)
	}
	fmt.Println("TCP server started and Listening on port", port)
	s.listener = listener

	s.wg.Add(1)
	go s.serve(db)
	return s
}

func (s *TcpServer) serve(db domain.KeyValueDB) {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				fmt.Printf("Failed to accept connection: %v\n", err)
			}
		} else {
			fmt.Println("Client connected")
			s.wg.Add(1)
			go func() {
				s.handleConnection(conn, db)
				s.wg.Done()
			}()
		}
	}
}

func (s *TcpServer) Stop() {
	close(s.shutdown)
	s.listener.Close()
	s.wg.Wait() // wait for active connections to complete
}

func (s *TcpServer) handleConnection(conn net.Conn, db domain.KeyValueDB) {
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

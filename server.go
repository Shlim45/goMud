package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Session struct {
	id   string
	conn net.Conn
}

func (s *Session) SessionId() string {
	return s.id
}

// TODO(jon): non-blocking write to session
func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

var nextSessionId = 1

func generateSessionId() string {
	var sid = nextSessionId
	nextSessionId++
	return fmt.Sprintf("%d", sid)
}

func handleConnection(conn net.Conn, inputChannel chan SessionEvent) error {
	log.Println("Connection accepted")

	buf := make([]byte, 4096)

	session := &Session{generateSessionId(), conn}

	inputChannel <- SessionEvent{session, &SessionCreatedEvent{}}

	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			inputChannel <- SessionEvent{session, &SessionDisconnectedEvent{}}
			log.Fatalln("Error reading from connection", err)
			return err
		}
		if n == 0 {
			log.Println("Closing connection")
			inputChannel <- SessionEvent{session, &SessionDisconnectedEvent{}}
			break
		}
		input := strings.TrimSpace(string(buf[0 : n-2]))
		log.Println("Received message:", input)

		inputChannel <- SessionEvent{session, &SessionInputEvent{input}}
	}

	return nil
}

func startServer(eventChannel chan SessionEvent) error {
	log.Println("Starting server")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// TODO(jon): handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// TODO(jon): handle error
			log.Fatalln("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
	//return nil
}

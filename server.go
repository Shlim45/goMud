package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var nextSessionId = 1

func generateSessionId() string {
	var sid = nextSessionId
	nextSessionId++
	return fmt.Sprintf("%d", sid)
}

func handleConnection(conn net.Conn, inputChannel chan SessionEvent) error {
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		ip := addr.IP.String()
		log.Println(fmt.Sprintf("Connection accepted from '%s'", ip))
	}

	buf := make([]byte, 4096)

	session := &Session{generateSessionId(), conn, DEFAULT}

	inputChannel <- SessionEvent{session, &SessionCreatedEvent{}}

	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			inputChannel <- SessionEvent{session, &SessionDisconnectedEvent{}}
			log.Printf("Session disconnected from %s", session.conn.LocalAddr())
		}
		if n == 0 {
			inputChannel <- SessionEvent{session, &SessionDisconnectedEvent{}}
			break
		}
		input := strings.TrimSpace(string(buf[0 : n-2]))

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
}

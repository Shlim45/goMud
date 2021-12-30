package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) error {
	log.Println("Connection accepted")

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("Connection ended by client", err)
			return err
		}
		if n == 0 {
			log.Println("Closing connection")
			break
		}
		msg := strings.TrimSpace(string(buf[0 : n-2]))
		log.Println("Received message:", msg)

		resp := fmt.Sprintf("You said \"%s\"\r\n", msg)
		n, err = conn.Write([]byte(resp))
		if err != nil {
			log.Fatalln("Error writing to connection", err)
			return err
		}
		if n == 0 {
			log.Println("Closing connection")
			break
		}
	}

	return nil
}

func startServer() error {
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
			if err := handleConnection(conn); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
	//return nil
}

func main() {
	err := startServer()
	if err != nil {
		log.Fatalln(err)
	}
}

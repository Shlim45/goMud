package main

import (
	"log"
	"net"
)

func handleConnection(conn net.Conn) error {
	log.Println("Connection accepted")
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

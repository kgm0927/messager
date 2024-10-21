package main

import (
	serverclient "github/messager/TCP/Server_client"
	"log"
	"net"
)

func main() {

	s := serverclient.NewServer()
	go s.Run(jsonserialize)

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		c := s.NewClient(conn)
		go c.ReadInput()
	}
}

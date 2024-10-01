package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	s := newServer()
	go s.run()
	fmt.Println("연결하고자 하는 주소를 입력해주세요.")
	var ipv4 string
	scanner := bufio.NewScanner(os.Stdin)

	// ipv4주소를 입력
	if scanner.Scan() {
		ipv4 = scanner.Text()
		ipv4 += ":8080"
		fmt.Println("입력한 주소:", scanner.Text())
	}

	listener, err := net.Listen("tcp", ipv4)
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

		c := s.newClient(conn)
		go c.readInput()
	}
}

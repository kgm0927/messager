package main

import (
	"bufio"
	"fmt"
	jsonfile "github/messager/TCP/json_file"
	"net"
	"os"
	"sync"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) // 스캐너 입력
	input := new(string)

	fmt.Println("연결할 ip를 입력하세요.")
	if scanner.Scan() { //입력받은 문자열 처리

		*input = scanner.Text()
		*input += ":8080"

	}

	listener, err := net.Listen("tcp", *input) // 네트워크 연결
	var jsm jsonfile.Making_message
	mtx := sync.Mutex{}

	if err != nil {
		fmt.Printf("%s\n", err)
	}
	defer listener.Close()

	done := make(chan struct{})

	for {

		conn, err := listener.Accept()
		if conn != nil {
			fmt.Println("연결 완료!")
		}
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		go handleConnection(conn, &mtx, scanner, jsm, done)
	}
}
func handleConnection(conn net.Conn, mtx *sync.Mutex, scanner *bufio.Scanner, jsm jsonfile.Making_message, done chan struct{}) {
	defer conn.Close()
	receive := make([]byte, 1024) // 수신할 데이터 버퍼

	// 데이터 수신 처리
	n, err := conn.Read(receive)
	if err != nil {
		fmt.Println("데이터 수신 오류:", err)
		return
	}

	// 수신된 데이터 출력
	mtx.Lock()
	fmt.Println("수신된 데이터:", string(receive[:n])) // 실제로 읽은 만큼만 출력
	mtx.Unlock()

	// 사용자 입력 처리
	if scanner.Scan() {
		mtx.Lock()
		input := scanner.Text()
		if input != "" {
			// JSON 메시지 직렬화 및 전송
			go func(c net.Conn, in string) {
				defer mtx.Lock()   // 데이터 전송 전에 잠금
				defer mtx.Unlock() // 데이터 전송 후 잠금 해제

				B := []byte(in) // 단순 문자열을 전송 (JSON 직렬화 생략)
				_, err := c.Write(B)
				if err != nil {
					fmt.Println("데이터 전송 오류:", err)
				}
			}(conn, input)
		}
		mtx.Unlock()
	}
}

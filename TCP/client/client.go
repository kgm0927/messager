package main

import (
	"bufio"
	"fmt"
	jsonfile "github/messager/TCP/json_file"
	"net"
	"os"
	"reflect"
	"sync"
)

func main() {

	var

	listener, err := net.Listen("tcp", "127.0.0.1:") // 네트워크 연결
	mtx := sync.Mutex{}
	if err != nil {
		fmt.Errorf("%s", err)
	}
	defer listener.Close()

	scanner := bufio.NewScanner(os.Stdin) // 스캐너 입력
	done := make(chan struct{})

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("%s", err)
			return
		}

	}
}
func handleConnection(conn net.Conn, mtx *sync.Mutex, scanner *bufio.Scanner, jsm jsonfile.Making_message, done chan struct{}) {
	defer conn.Close()

	var input string
	if scanner.Scan() { //입력받은 문자열 처리

		input = scanner.Text()
	}

	defer mtx.Unlock() // 작업이 끝나면 반드시 잠금 해제

	if input != "" { // 문자가 오면 잠그기
		mtx.Lock()
	}
	// json 입력 및 작성
	go func() {

		send := jsonfile.Message{Talk: input, Room_name: "", Client_name: ""} //메시지 구조체 생성
		jsm = &send
		err, B := jsm.Serialize()

		if err != nil {
			fmt.Print(err, "\n")
		}
		conn.Write(B)

	}()

	// json 해제 및 읽기
	mtx.Lock()
	go func() {
		receive := make([]byte, reflect.TypeOf(jsonfile.Message{}).Size())
		_, err := conn.Read(receive)
		if err != nil {
			fmt.Print(err, "\n")
		}

		err, message := jsm.UnSerialize(receive)
		fmt.Println(message.Talk)
	}()
}

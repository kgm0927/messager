package main

import (
	"bufio"
	"fmt"
	jsonfile "github/messager/TCP/json_file"
	"net"
	"os"
	"reflect"
	"strings"
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
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		go handleConnection(conn, &mtx, scanner, jsm, done)
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
	mtx.Unlock()

	// json 해제 및 읽기
	mtx.Lock()
	go func() {
		receive := make([]byte, reflect.TypeOf(jsonfile.Message{}).Size())
		_, err := conn.Read(receive)
		if err != nil {
			fmt.Print(err, "\n")
		}

		err, message := jsm.UnSerialize(receive)
		if err != nil {
			fmt.Printf(" 알아들을 수 없음: %s\n", err)
		}
		fmt.Println(message.Talk)

		if strings.Contains(message.Talk, "/quit") {
			fmt.Println("프로그램을 종료합니다.")
			close(done)
		}
	}()
	mtx.Unlock()
}

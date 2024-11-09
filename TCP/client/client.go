package main

import (
	"bufio"
	"bytes"
	"fmt"

	"log"
	"net"
	"os"
	"strings"
	"sync"

	serverclient "github.com/messager/TCP/Server_client"
	jsonfile "github.com/messager/TCP/json_file"
)

func main() {
	// 서버에 연결할 IP와 포트 설정

	scanner := bufio.NewScanner(os.Stdin) // 사용자 입력을 위한 스캐너 생성
	var address string
	basic_json_file := new(jsonfile.Message)
	if scanner.Scan() {
		address += scanner.Text()
		address += ":8080"
	}

	conn, err := net.Dial("tcp", address) // TCP 연결
	if err != nil {
		fmt.Printf("서버에 연결할 수 없습니다: %s\n", err)
		return
	}
	defer conn.Close() // 클라이언트 종료 시 연결 닫기

	basic_json_file.Room_name = "nothing"
	basic_json_file.Client_name = "anonymous"
	basic_json_file.Ip = conn.LocalAddr().String()

	var wg sync.WaitGroup
	done := make(chan struct{})

	fmt.Println("서버에 연결되었습니다. 메시지를 입력하세요.")

	for {
		wg.Add(2)
		go readResponses(&conn, &wg) // 서버 응답을 읽는 고루틴 시작
		go WriteResponse(scanner, &conn, &wg, done, basic_json_file)
	}
	if _, boolean := <-done; boolean {
		fmt.Println("프로그램을 종료하겠습니다.")

		return
	}

	wg.Wait()
	// 모든 고루틴이 종료된 후 done 채널을 닫음

}

func WriteResponse(scanner *bufio.Scanner, conn *net.Conn, wg *sync.WaitGroup, done chan struct{}, mssg *jsonfile.Message) {
	defer wg.Done()

	for {
		var err error
		scanner.Scan()          // 사용자 입력 대기
		input := scanner.Text() // 입력할 문자 출력

		if input == "exit" { // 'exit' 입력 시 종료
			fmt.Println("클라이언트 종료.")
			fmt.Fprintf(*conn, "%s\n", input)
			done <- struct{}{}
			return

		}

		sentences := strings.Split(input, " ")

		mssg.Args = strings.Join(sentences[1:], " ")
		mssg.Id, err = serverclient.Compare_cmd(sentences[0])

		fmt.Println(mssg.Args, mssg.Id, mssg.Ip, mssg.Room_name)

		if err != nil {
			fmt.Println(err)
		}

		mb, err := mssg.Serialize()
		if err != nil {
			fmt.Println(err)
		}

		// 메시지를 서버로 전송
		n, err := (*conn).Write(mb)
		if err != nil {
			fmt.Printf("메세지를 서버로 전송 불가: %s", err)
		} else {
			fmt.Println(input, "전송완료, size:", n)
		}
	}
}

// 서버로부터의 응답을 읽는 함수
func readResponses(conn *net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(*conn)
	if conn == nil {
		log.Fatal("Connection is nil") // 연결이 nil일 경우 종료
	}

	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println("Failed to read from connection:", err)
			break
		}
		line = bytes.Trim(line, "\x00") // 필요 시 trim
		// JSON 역직렬화
		msg := new(jsonfile.Message)
		if msg, err = msg.UnSerialize(line); err != nil {
			log.Println("UnSerialize error:", err)
			continue
		}
		fmt.Println("서버로부터 수신된 메시지:", msg.Args)
	}
}

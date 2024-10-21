package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	// 서버에 연결할 IP와 포트 설정

	scanner := bufio.NewScanner(os.Stdin) // 사용자 입력을 위한 스캐너 생성
	var address string

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

	var wg sync.WaitGroup
	done := make(chan struct{})

	fmt.Println("서버에 연결되었습니다. 메시지를 입력하세요.")
	wg.Add(2)

	go WriteResponse(scanner, &conn, &wg, done)
	go readResponses(conn, &wg) // 서버 응답을 읽는 고루틴 시작

	if _, boolean := <-done; boolean {
		fmt.Println("프로그램을 종료하겠습니다.")

		return
	}

	wg.Wait()
	// 모든 고루틴이 종료된 후 done 채널을 닫음

}

func WriteResponse(scanner *bufio.Scanner, conn *net.Conn, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()

	for {

		scanner.Scan() // 사용자 입력 대기
		input := scanner.Text()

		if input == "exit" { // 'exit' 입력 시 종료
			fmt.Println("클라이언트 종료.")
			fmt.Fprintf(*conn, "%s\n", input)
			done <- struct{}{}
			return

		}

		// 메시지를 서버로 전송
		_, err := fmt.Fprintf(*conn, "%s\n", input)
		if err != nil {
			continue
		} else {
			fmt.Println(input)
		}
	}
}

// 서버로부터의 응답을 읽는 함수
func readResponses(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		reader := bufio.NewReader(conn) // 수신 리더 생성

		response, err := reader.ReadString('\n') // 서버로부터의 응답 수신
		if err != nil {
			fmt.Println("서버와의 연결이 종료되었습니다:", err)
			break
		}
		fmt.Println("서버로부터 수신된 메시지:", response) // 수신된 메시지 출력

	}
}

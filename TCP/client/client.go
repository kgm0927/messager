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
	address := "localhost:8080"
	var wg sync.WaitGroup
	conn, err := net.Dial("tcp", address) // TCP 연결
	if err != nil {
		fmt.Printf("서버에 연결할 수 없습니다: %s\n", err)
		return
	}
	done := make(chan struct{})

	defer conn.Close() // 클라이언트 종료 시 연결 닫기

	fmt.Println("서버에 연결되었습니다. 메시지를 입력하세요.")

	scanner := bufio.NewScanner(os.Stdin) // 사용자 입력을 위한 스캐너 생성

	wg.Add(2)
	go WriteResponse(scanner, &conn, &wg, done)
	go readResponses(conn, &wg) // 서버 응답을 읽는 고루틴 시작
	if _, boolean := <-done; boolean {
		return
	}
	wg.Wait()

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
			continue

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
	reader := bufio.NewReader(conn) // 수신 리더 생성
	for {
		response, err := reader.ReadString('\n') // 서버로부터의 응답 수신
		if err != nil {
			fmt.Println("서버와의 연결이 종료되었습니다:", err)
			return
		}
		fmt.Println("서버로부터 수신된 메시지:", response) // 수신된 메시지 출력
	}
}


네트워크와 그에 대한 패킷인 TCP를 제대로 이해하기 위해 메신저를 개발해 보기로 결정하였다. 다음과 같은 환경에서 네트워크를 통해 서로 문자를 주고받는 시스템을 개발할 계획이다.
그리고 부가적으로 Go 안에 있는 goroutine과 channel 키워드에 대해서도 짧게나마 다뤄볼 계획이다.

### 환경

- 프로그래밍 언어: Go
- 운영체제 윈도우: Windwows
- IDE: Visual Studio Code
- 그외 응용 프로그램: PuTTY


---
# 작동 원리

응용 프로그램 PuTTY에서는 telnet라는 원격지의 호스트 컴퓨터에 접속하기 위해 사용하는 인터넷 프로토콜을 통해 통신할 계획이다. 아래와 같은 커맨드를 만들어 상태편과 대화를 시도할 계획이다.


#### 커맨드

- `/nick [name]`- 이름을 정하게 됨. 
- `/join [name]`- 들어갈 채팅방의 이름을 정하게 됨, 만약 방이 없다면, 새로 만들어야 함. 유저는 오로지 한 방만 들어갈 수 있다.
- `/rooms`-모든 방의 리스트를 확인할 수 있음
- `/msg <msg>`- 방에 있는 모든 이들에게 메시지(broadcast message)를 보낼 수 있음.
- `/quit`- 채팅 서버와 연결을 끊음


---
# 코드 

코드의 내용 순서는 이러한 방식으로 흘러갈 것이다.

- room.go: 채팅 방을 구현하는 코드
- client.go: 현재 유저와 연결된 방을 확인할 수 있음.
- command.go: 명령어를 정의 및 정형화를 담당함.
- server.go: 클라이언트로부터 들어오는 명령어를 인식하고 서버 내 클라이언트를 담당함.
- main.go: 최종적으로 프로그램을 실행하는 역할



### room.go

``` go
package main

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}
```

- type room struct

채팅 방을 정의하는 구조체이다. 변수를 보면 문자열 타입의 **name**, 해시 맵 배열인 members가 있다. 해시 맵 members는 서버라는 채팅방에 접속하고 있는 클라이언트와 그의 ip주소를 저장해 놓는 곳이다.

- `broadcast(sender *client, msg string)`

클라이언트가 보낸 메시지를 받아들이고 그 외 모든 클라이언트(채팅방의 멤버)들에게 보내는 역할을 함. 만약 보내는 이와 서버의 ip주소가 다르면 메시지를 보내는데, 이는 문자를 보낸 클라이언트에게 다시 문자를 보낼 필요가 없기 때문이다.



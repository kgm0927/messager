
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





### Testing.go


``` go
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
```


- `type client struct`

클라이언트를 정의하는 구조체이다. 변수를 하나씩 분석해보자.

1. `conn net.Conn`

conn은 net 패키지의 Conn 인터페이스를 정의한다. Conn 인터페이스는 원격(remote) 및 지역(local) 네트워크와 연결하여 문자열을 주고 받는 역할을 한다.

2. `nick string`

클라이언트의 이름을 저장하는 변수이다.


3. `room *room`

파일 room.go에 room 구조체이다. 이를 통해 자신이 들어가고자 하는 채팅방의 주소를 저장해 놓으며, 한번에 하나의 채팅방만 들어갈 수 있다.

4. `commands chan<- command`

커맨드를 보내는 역할을 하는 채널 변수이다. 오로지 보내는 것만 가능하며, 받는 기능을 수행하지는 못한다. 이 commands 변수는 나중에 room에서 인자로 받아들여, broadcast 함수를 실행할 수 있도록 한다.


- `func (c *client) readInput() `

readInput() 함수는 거대한 for문으로 계속 서버 쪽에서 보내는 문자를 읽어오는 역할을 한다.
바로 `msg, err := bufio.NewReader(c.conn).ReadString('\n')`가 서버로부터 문자열을 읽어온다.

문자열 변수 msg를 `strings.Trim('\r\n')` 함수를 통해 개행된 문자열을 제거하고 msg을 명령어와 그외의 문자열을 분석하여 어떤 명령을 내릴 지 switch 문을 통해 판단한다.


- `func (c *client) err(err error)`

클라이언트에서 처리하지 못하면 에러를 보내는 역할을 한다. `c.conn.Write([]byte("err: " + err.Error() + "\n"))`를 통해 생성된 에러가 서버로 전달된다.


- `func (c *client) msg(msg string)`
문자를 보내는 역할을 한다. `c.conn.Write([]byte("> " + msg + "\n"))` 통해 문자를 전달한다.


### command.go


``` go
package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
```

커맨드 명령을 정형화하는 역할을 한다. 커맨드는 그냥 문자가 아니라 상수이며, 이에 맞춰서 Client.go의 ReadInput()에서 커맨드를 분석하면 어떤 행위를 해야 할지 결정한다.

- CMD_NICK
'/nick'을 입력한 후 뒤에 이름을 붙여서 어떤 닉네임을 만들 것인지 결정한다.

- CMD_JOIN
'/join' 명령어 다음에 이름을 붙여서 어디에 들어갈 지 결정한다. 만약 그 방이 존재하지 않는다면 새로 만들어야 한다.

- CMD_ROOMS
'/rooms'를 통해 현재 서버에 방이 몇개나 있는지 확인해 볼 수 있다.

- CMD_MSG
'/msg' 명령어 다음에 전달하고자 하는 문자열을 생성해서 전달하는 역할을 한다. 서버로 성공적으로 전달이 되면 room.go 파일에 있는 broadcast() 함수를 통해 모든 서버에 있는 클라이언트들에게 메시지를 전달하는 역할을 한다.

- CMD_QUIT
'/quit' 명령어를 통해 현재 연결된 서버와의 통신을 끊는 역할을 한다.


- `type command struct`

커맨드 명령을 정의하는 역할을 한다. 그냥 명령만 전달하는 것이 아니라, 보내는 클라이언트가 누구인지, 그리고 커맨드 명령어 다음으로 오는 문자열을 분석해서 그대로 메시지로 보낼 지 아니면, 다른 명령을 수행할 지 결정한다.




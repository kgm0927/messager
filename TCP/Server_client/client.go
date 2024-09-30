package serverclient

import (
	"encoding/json"
	"fmt"
	jsonfile "github/messager/TCP/json_file"
	"io"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
	message  jsonfile.Message
}

func (c *client) Insert_message(M *jsonfile.Message) {
	c.message = *M
}

func (c *client) Return_Message() jsonfile.Message {
	return c.message
}

func (c *client) ReadInput(mi jsonfile.Making_message) { // 입력을 읽음.
	// mi는 인터페이스이므로, 구조체를 직접 할당할 수 없습니다.

	msg_client := new(jsonfile.Message)
	mi = msg_client

	B, err := io.ReadAll(c.conn)
	if err != nil {
		fmt.Println("메시지에 오류")
	}

	err, msg := (mi).UnSerialize(B)
	if err != nil {
		return
	}

	msg.Talk = strings.Trim(msg.Talk, "\r\n")

	args := strings.Split(msg.Talk, " ")
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
		(mi).Initialize(c.room.name, c.nick, args[1])
		c.Insert_message((mi).Return_self())

	case "/quit":
		c.commands <- command{
			id:     CMD_QUIT,
			client: c,
		}
	default:
		c.err(fmt.Errorf("unknown command: %s", cmd), mi)
	}
}

func (c *client) err(err error, Mi jsonfile.Making_message) {

	err, B := (Mi).Serialize()
	if err != nil {
		fmt.Print("json serialize 불가능")
		return
	}

	json_connect := json.NewDecoder(c.conn)
	if err := (*json_connect).Decode(&B); err != nil { // Decode함수를 쓸 때 파라미터는 반드시 call by address여야한다.
		fmt.Println("Error decoding JSON:", err)
		return
	}
}

func (c *client) msg(msg string, Mi *jsonfile.Making_message) { // 수정할 필요가 있어보임.
	c.message.Talk = msg
	*Mi = &c.message

	err, B := (*Mi).Serialize()
	if err != nil {
		fmt.Print("json serialize 불가능")
		return
	}

	json_connect := json.NewDecoder(c.conn)
	if err = (*json_connect).Decode(&B); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

}

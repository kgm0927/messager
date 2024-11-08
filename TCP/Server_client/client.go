package serverclient

import (
	"encoding/json"
	"fmt"
	"slices"

	"net"
	"strings"

	jsonfile "github.com/messager/TCP/json_file"
)

type Client struct {
	conn     net.Conn
	nick     string
	room     *room
	Commands chan command
}

// func (c *Client) Insert_message(M *jsonfile.Message) {
// 	c.Commands <- *M
// }

// func (c *Client) Return_Message() chan<- jsonfile.Message {
// 	return c.Commands
// }

func (c *Client) ReadInput() { // 입력을 읽음.

	for { // mi는 인터페이스이므로, 구조체를 직접 할당할 수 없습니다.
		var reading jsonfile.Message
		fmt.Println("읽고 있는 중 ,ReadInput")
		var err error
		comm := new(command)

		err = json.NewDecoder(c.conn).Decode(&reading) // 네크워크로 수신한 코드 복호화(decode)
		if err != nil {
			fmt.Println("수신불가")
		}

		reading.Args = strings.TrimSpace(reading.Args) // 쓸데없는 것 지우기
		comm.args = reading.Args

		slice := strings.Split(comm.args, " ") // /(명령어) 입력을 위해 잠시 쪼갬

		cmd := reading.Id // 메시지 종료 알려줌
		comm.id = commandID(cmd)
		comm.client = c

		switch cmd {
		case int(CMD_NICK): // /nick
			slice = slices.Insert[[]string, string](slice, 0, "/nick")
			comm.args = strings.Join(slice, " ")

			c.Commands <- *comm

		case int(CMD_JOIN): // /join
			slice = slices.Insert[[]string, string](slice, 0, "/join")
			comm.args = strings.Join(slice, " ")

			c.Commands <- *comm

		case int(CMD_ROOMS): // /rooms
			slice = slices.Insert[[]string, string](slice, 0, "/rooms")
			comm.args = strings.Join(slice, " ")

			c.Commands <- *comm

		case int(CMD_MSG): // /msg
			slice = slices.Insert[[]string, string](slice, 0, "/msg")
			comm.args = strings.Join(slice, " ")

			c.Commands <- *comm

		case int(CMD_QUIT): // /quit
			slice = slices.Insert[[]string, string](slice, 0, "/quit")
			comm.args = strings.Join(slice, " ")

			c.Commands <- *comm

		default:

			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *Client) err(err error) {
	// 다시 작성할 것
}
func (c *Client) Setting_message_sentence(msg string) (cmd string, args []string) {
	msg = strings.Trim(msg, "\n") // 수정할 가능성이 있음.

	args = strings.Split(msg, " ")
	cmd = strings.TrimSpace(args[0])

	var msg_new []string
	if len(msg) > 1 {
		msg_new = append(args[1:], args[2:]...)
	} else {
		msg_new = []string{cmd} // 인자가 하나만 있는 경우
	}

	return cmd, msg_new
}

func (c *Client) Setting_Message_file(cmd string, msg []string) jsonfile.Message { // 사용하지 않음
	setting := new(jsonfile.Message)

	cmd_num, err := Compare_cmd(cmd)

	if err != nil {
		fmt.Printf("%s", err)
	}

	setting.Id = cmd_num
	setting.Args = strings.Join(msg, " ") // 나중에 문자열에 '>' 붙임
	setting.Client_name = c.nick
	setting.Room_name = c.room.name
	setting.Ip = c.conn.LocalAddr().String()

	return *setting

}

func (c *Client) msg(msg string) { // 수정할 필요가 있어보임.

	// ----------------------------------------------- 문자열 분석 및 id 출력

	str := strings.Split(msg, " ")

	sending := c.Setting_Message_file(str[0], str[1:])
	B, err := sending.Serialize()

	if err != nil {
		fmt.Println("%s", err)
	}

	c.conn.Write(B)

}

func Compare_cmd(cmd string) (int, error) {

	switch cmd {
	case "/nick": //0
		fmt.Println("전달은 0, /nick")
		return int(CMD_NICK), nil

	case "/join": // 1
		fmt.Println("전달은 1, /join")

		return int(CMD_JOIN), nil

	case "/rooms": // 2
		fmt.Println("전달은 2, /rooms")

		return int(CMD_ROOMS), nil

	case "/msg": // 3
		fmt.Println("전달은 3, /msg")

		return int(CMD_MSG), nil

	case "/quit": // 4
		fmt.Println("전달은 4, /quit")

		return int(CMD_QUIT), nil

	default:
		return 100, fmt.Errorf("올바른 문자가 아님")
	}
}

func (c *Client) Return_connent() *net.Conn {
	return &c.conn
}

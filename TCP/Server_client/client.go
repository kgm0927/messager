package serverclient

import (
	"encoding/json"
	"fmt"
	jsonfile "github/messager/TCP/json_file"

	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	Commands chan jsonfile.Message
}

func (c *client) Insert_message(M *jsonfile.Message) {
	c.Commands <- *M
}

func (c *client) Return_Message() chan<- jsonfile.Message {
	return c.Commands
}

func (c *client) ReadInput() { // 입력을 읽음.
	// mi는 인터페이스이므로, 구조체를 직접 할당할 수 없습니다.
	var reading jsonfile.Message
	err := json.NewDecoder(c.conn).Decode(&reading)

	if err != nil {
		fmt.Println(err)
	}
	reading.Args = strings.Trim(reading.Args, "/n")
	cmd := reading.Id

	switch cmd {
	case int(CMD_NICK):
		c.Commands <- reading

	case int(CMD_JOIN):
		c.Commands <- reading

	case int(CMD_ROOMS):
		c.Commands <- reading

	case int(CMD_MSG):
		c.Commands <- reading

	case int(CMD_QUIT):
		c.Commands <- reading
	default:
		c.err(fmt.Errorf("unknown command: %s", cmd))
	}
}

func (c *client) err(err error) {
	// 다시 작성할 것
}
func (c *client) Setting_message_sentence(msg string) (cmd string, args []string) {
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

func (c *client) Setting_Message_file(cmd string, msg []string) jsonfile.Message {
	setting := new(jsonfile.Message)

	cmd_num, err := Compare_cmd(cmd)

	if err != nil {
		fmt.Printf("%s", err)
	}

	setting.Id = cmd_num
	setting.Args = strings.Join(msg, " ") // 나중에 문자열에 '>' 붙임
	setting.Client_name = c.nick
	setting.Room_name = c.room.name

	return *setting

}

func (c *client) msg(msg string) { // 수정할 필요가 있어보임.

	cmd, args := c.Setting_message_sentence(msg)

	// ----------------------------------------------- 문자열 분석 및 id 출력

	sending := c.Setting_Message_file(cmd, args)
	B, err := sending.Serialize()

	if err != nil {
		fmt.Println("%s", err)
	}

	c.conn.Write(B)

}

func Compare_cmd(cmd string) (int, error) {

	switch cmd {
	case "/nick": //0

		return int(CMD_NICK), nil

	case "/join": // 1

		return int(CMD_JOIN), nil

	case "/rooms": // 2

		return int(CMD_ROOMS), nil

	case "/msg": // 3

		return int(CMD_MSG), nil

	case "/quit": // 4

		return int(CMD_QUIT), nil

	default:
		return 100, fmt.Errorf("올바른 문자가 아님")
	}
}

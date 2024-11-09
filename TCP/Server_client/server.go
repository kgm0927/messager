package serverclient

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms map[string]*room
	// buffer_client []*Client
	commands chan command
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

// func (s *server) searching_member_client(member *jsonfile.Message) (*Client, error) {
// 	room, exists := s.rooms[member.Room_name]
// 	if !exists {
// 		return nil, fmt.Errorf("방 %s가 존재하지 않습니다", member.Room_name)
// 	}

// 	for _, client := range room.members {
// 		if client.conn.RemoteAddr().String() == member.Ip {
// 			return client, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("클라이언트 %s가 방 %s에 없습니다", member.Ip, member.Room_name)
// }

// func (s *server) find_client_no_room(cmd *jsonfile.Message) *Client {
// 	for i, v := range s.buffer_client {
// 		if v.conn.LocalAddr().String() == cmd.Ip {

// 			client := s.buffer_client[i]

// 			return client
// 		}
// 	}
// 	return nil
// }

func (s *server) Run() { // 이것이 계속 돌아감
	for cmd := range s.commands {

		switch int(cmd.id) {
		case int(CMD_NICK):
			fmt.Println("닉네임 만들기")

			s.nick(cmd.client, cmd.args) // 메시지 전달 //

		case int(CMD_JOIN):
			s.join(cmd.client, cmd.args) // 방 입장 //

		case int(CMD_ROOMS):
			s.listRooms(cmd.client) // //

		case int(CMD_MSG):
			fmt.Println("메시지 실행")
			s.msg(cmd.client, cmd.args) //

		case int(CMD_QUIT):
			s.quit(cmd.client) //
		}
	}
}

func (s *server) NewClient(conn net.Conn) *Client { // 새로운 클라이언트를 만듦 // 이것은 1번만 실행됨.

	log.Printf("new Client has joined: %s", (conn).RemoteAddr().String())

	c := &Client{ // 클라이언트 생성
		conn:     (conn),
		nick:     "anonymous",
		room:     nil,
		Commands: s.commands,
	}

	/////////////////////  메세지 생성

	// s.buffer_client = append(s.buffer_client, c)

	return c
}

func (s *server) nick(c *Client, cmd string) {
	if len(cmd) < 2 || cmd == "anonymous" {
		c.msg("nick is required. usage: /nick NAME") // 메시지 보냄
		return
	}
	fmt.Println("닉네임 만들기 시작")
	c.nick = cmd
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *Client, cmd string) {

	if c.nick == "" || c.nick == "anonymous" {
		fmt.Println("이름을 설정해 주시기 바랍니다. ")
		return
	}

	if c.room.name == "nothing" || cmd == "nothing" {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := cmd

	r, ok := s.rooms[roomName]
	if !ok {

		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*Client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *Client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *Client, msg string) {

	if c.nick == "" || c.nick == "anonymous" {
		fmt.Println("이름을 설정해 주시기 바랍니다. ")
		return
	}

	if c.room == nil {
		fmt.Println("들어갈 방을 먼저 선택하세요.")
		return
	}

	if len(msg) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	// msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg) // 메시지 수정 전달 // room.go로 이동
}

func (s *server) quit(c *Client) {
	log.Printf("Client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *Client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

package serverclient

import (
	"fmt"
	jsonfile "github/messager/TCP/json_file"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan jsonfile.Message
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan jsonfile.Message),
	}
}
func (s *server) searching_member_client(member *jsonfile.Message) *client {
	room := s.rooms[member.Room_name]

	for address, client := range room.members {
		if address.String() == member.Ip {
			return client
		} else {
			return nil
		}
	}
	return nil
}

func (s *server) Run() {
	for cmd := range s.commands { // 여기서 받은 json파일을 분해해야만 한다.

		client := s.searching_member_client(&cmd)
		msg := cmd.Args
		switch cmd.Id {
		case int(CMD_NICK):
			s.nick(client, cmd) // 메시지 전달

		case int(CMD_JOIN):
			s.join(client, cmd) // 방 입장

		case int(CMD_ROOMS):
			s.listRooms(client) //

		case int(CMD_MSG):
			s.msg(client, msg)

		case int(CMD_QUIT):
			s.quit(client)
		}
	}
}

func (s *server) NewClient(conn net.Conn) *client {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		room:     nil,
		Commands: s.commands,
	}
}

func (s *server) nick(c *client, cmd jsonfile.Message) {
	if len(cmd.Client_name) < 2 || cmd.Client_name == "anonymous" {
		c.msg("nick is required. usage: /nick NAME") // 메시지 보냄
		return
	}

	c.nick = cmd.Args
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, cmd jsonfile.Message) {
	if len(cmd.Args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := cmd.Args

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, msg string) {

	if c.room == nil {
		fmt.Println("들어갈 방을 먼저 선택하세요.")
	}

	if c.nick == "" || c.nick == "anonymous" {
		fmt.Println("이름을 설정해 주시기 바랍니다. ")
	}

	if len(msg) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	// msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg) // 메시지 수정 전달 // room.go로 이동
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

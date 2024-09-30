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
	commands chan command
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Run(jfm *jsonfile.Making_message) {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args, jfm)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args, jfm)
		case CMD_ROOMS:
			s.listRooms(cmd.client, jfm)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args, jfm)
		case CMD_QUIT:
			s.quit(cmd.client, jfm)
		}
	}
}

func (s *server) NewClient(conn net.Conn) *client {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
}

func (s *server) nick(c *client, args []string, jfm *jsonfile.Making_message) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME", jfm) // 메시지 보냄
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick), jfm)
}

func (s *server) join(c *client, args []string, jfm *jsonfile.Making_message) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME", jfm)
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c, jfm)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick), jfm)

	c.msg(fmt.Sprintf("welcome to %s", roomName), jfm)
}

func (s *server) listRooms(c *client, jfm *jsonfile.Making_message) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")), jfm)
}

func (s *server) msg(c *client, args []string, jfm *jsonfile.Making_message) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG", jfm)
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg, jfm)
}

func (s *server) quit(c *client, jfm *jsonfile.Making_message) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c, jfm)

	c.msg("sad to see you go =(", jfm)
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client, jfm *jsonfile.Making_message) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick), jfm)
	}
}

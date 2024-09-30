package serverclient

import (
	jsonfile "github/messager/TCP/json_file"
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string, jfm *jsonfile.Making_message) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg, jfm)
		}
	}
}

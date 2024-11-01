package serverclient

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*Client
}

func (r *room) broadcast(sender *Client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg) // client.go 이동
		}
	}
}

package serverclient

import (
	"io"
	"net"
	"testing"
)

func Testclient(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		t.Fatal(err)
	}

	C := client{conn:&net.TCPConn{},"김민",room: &room{name:"1",},commands: "메시지" ,message:}
	var jmm jsonfile.Making_message

	C.Insert_message(&M)
	
	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024)

				for {
					n, err := c.Read(buf) // #4
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

}

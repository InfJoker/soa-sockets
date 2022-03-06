package main

import (
	"fmt"
	"net"
	"io"
)


var PORT = ":5454"


func handler(c net.Conn, m map[string]net.Conn) {
	var login string
	fmt.Fscanln(c, &login)
	m[login] = c

	for {
		var bytesNum int
		n, err := fmt.Fscanln(c, &bytesNum)
		if n == 0 || err != nil {
			if err.Error() == "EOF" {
				return
			}
			continue
		}
		audioBytes := make([]byte, bytesNum)

		io.ReadFull(c, audioBytes)

		go broadcast(m, login, audioBytes)
	}
}


func broadcast(m map[string]net.Conn, login string, audioBytes []byte) {
	for receiver_login, conn := range m {
		if receiver_login != login {
			fmt.Fprintln(conn, len(audioBytes))
			conn.Write(audioBytes)
			fmt.Printf("Send %d bytes to %s\n", len(audioBytes), login)
		}
	}
}


func main() {
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	m := make(map[string]net.Conn)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handler(c, m)
	}
}

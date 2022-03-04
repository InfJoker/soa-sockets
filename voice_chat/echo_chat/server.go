package main

import (
	"fmt"
	"net"
	"io"
)


var PORT = ":5454"


func handler(c net.Conn) {
	var login string
	fmt.Fscanln(c, &login)

	action := "start"

	for action != "quit" {
		var bytesNum int
		fmt.Fscanln(c, &bytesNum)

		audioBytes := make([]byte, bytesNum)
		io.ReadFull(c, audioBytes)

		fmt.Fprintln(c, len(audioBytes))

		c.Write(audioBytes)

		fmt.Fscanln(c, &action)
	}
}


func main() {
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handler(c)
	}
}

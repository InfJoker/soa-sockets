package main

import (
	"fmt"
	"net"
)

var PORT = ":5454"

func handler(c net.Conn) {
	fmt.Fprintln(c, "Hello World!!!")
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

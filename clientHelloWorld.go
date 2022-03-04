package main

import (
	"fmt"
	"net"
)

var SERVER_ADDRESS = "127.0.0.1:5454"

func main() {
	c, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	b := make([]byte, 64)
	c.Read(b)
	fmt.Println(string(b))
}

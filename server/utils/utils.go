package utils

import (
	"fmt"
	"net"
)

type clientHandle struct {
	Id   string
	Conn net.Conn
}

type servTy struct {
	L       net.Listener
	Clients map[string]clientHandle
}

var server *servTy

func handleClient(c net.Conn) {
}

func Run(portNr string) {
	server = &servTy{}

	portStr := ":" + portNr
	l, err := net.Listen("tcp4", portStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	server.L = l

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleClient(c)
	}
}

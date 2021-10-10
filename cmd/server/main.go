package server

import (
	"fmt"
	"net"
)

func handleClient(c net.Conn) {

}

func Run(portNr string) {
	portStr := ":" + portNr
	l, err := net.Listen("tcp4", portStr)
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
		go handleClient(c)
	}
}

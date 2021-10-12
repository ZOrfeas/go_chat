package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	client "github.com/ZOrfeas/go_chat/client/utils"
	common "github.com/ZOrfeas/go_chat/common/utils"
)

type servTy struct {
	L       net.Listener
	Clients map[string]*client.CliTy
	Ctx     context.Context
	Cancel  context.CancelFunc
}

var server *servTy

func handleStdinInput(stdin <-chan string) {
	for cmd := range stdin {
		res, err := handleServerCommand(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res)
	}
}
func setupClient(c net.Conn) (*client.CliTy, error) {
	log.Println("Setting up client")
	reader := bufio.NewReader(c)
	firstMessage, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(firstMessage, client.FirstMessagePrefix) {
		return nil, fmt.Errorf("internal: first client message was not '" +
			client.FirstMessagePrefix + "'")
	}
	words := strings.Fields(firstMessage)
	if len(words) != 2 {
		return nil, fmt.Errorf("internal: first client message doesn't have 2 fields")
	}
	candidateId := words[1]
	for _, idExists := server.Clients[candidateId]; idExists; {
		candidateId += "_"
	}
	return &client.CliTy{Id: candidateId, Conn: c}, nil
}
func teardownClient(thisClient *client.CliTy) {
	log.Println("Tearing down client", thisClient.Id)
	delete(server.Clients, thisClient.Id)
}

func handleClient(c net.Conn) {
	defer c.Close()
	log.Println("Incoming client:", c.RemoteAddr().String())

	thisClient, err := setupClient(c)
	if err != nil {
		log.Println(err)
		return
	}
	defer teardownClient(thisClient)

	clientIn := make(chan string, 1)
	defer close(clientIn)

	go common.ChannelStrings(clientIn, thisClient.Conn)
	for {
		select {
		case req := <-clientIn:
			{
				exit, err := handleClientRequest(thisClient, req)
				if err != nil {
					log.Println(err)
					continue
				}
				if exit {
					return
				}
			}
		case <-server.Ctx.Done():
			return // will close the client handle
		}
	}
}

func Run(portNr string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Attempting to listen on port", portNr)
	l, err := net.Listen("tcp4", ":"+portNr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	log.Println("Server listening on port", portNr)

	customCancel := func() { cancel(); l.Close() }
	server = &servTy{
		L: l, Ctx: ctx, Cancel: customCancel,
		Clients: map[string]*client.CliTy{},
	}

	stdin := make(chan string, 1)
	defer close(stdin)

	log.Println("Setting up command prompt processing")
	go common.ChannelStrings(stdin, os.Stdin)
	go handleStdinInput(stdin)
	log.Println("Command prompt processing successfuly setup")

	log.Println("Accepting incoming TCP connections")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleClient(c)
	}
}

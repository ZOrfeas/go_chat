package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	client "github.com/ZOrfeas/go_chat/client/utils"
	common "github.com/ZOrfeas/go_chat/common/utils"
)

type clientWrapperTy struct {
	Client *client.CliTy
	mu     sync.Mutex
}

func (thisClient *clientWrapperTy) checkAndSend(rep string) error {
	thisClient.mu.Lock()
	defer thisClient.mu.Unlock()
	rep = strings.ReplaceAll(rep, common.HostCommandIdentifier, "[redacted]")
	rep = strings.TrimRightFunc(rep, unicode.IsSpace)
	return thisClient.Client.SendString(rep)
}

func (thisClient *clientWrapperTy) sendCommand(cmd common.HostCommand, arg string) error {
	thisClient.mu.Lock()
	defer thisClient.mu.Unlock()
	rep := fmt.Sprint(
		common.HostCommandIdentifier+strconv.Itoa(int(cmd)),
		" ",
		arg,
	)
	return thisClient.Client.SendString(rep)
}

type servTy struct {
	L       net.Listener
	Clients map[string]*clientWrapperTy
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
func setupClient(c net.Conn) (*clientWrapperTy, error) {
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

	for cliWrapper, idExists := server.Clients[candidateId]; idExists; {
		if cliWrapper == nil {
			break
		}
		candidateId += "_"
	}
	newClient := &client.CliTy{Id: candidateId, Conn: c}
	newClientWrapper := &clientWrapperTy{Client: newClient}
	server.Clients[candidateId] = newClientWrapper
	newClient.SendString(candidateId)
	return newClientWrapper, nil
}
func teardownClient(thisClient *clientWrapperTy) {
	log.Println("Tearing down client", thisClient.Client.Id)
	delete(server.Clients, thisClient.Client.Id)
}

func handleClient(c net.Conn) {
	defer log.Println("Closing client", c.RemoteAddr().String())
	defer c.Close()
	log.Println("Incoming client:", c.RemoteAddr().String())

	thisClient, err := setupClient(c)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client", "'"+thisClient.Client.Id+"'", "setup")
	defer teardownClient(thisClient)

	clientIn := make(chan string, 1)
	go common.ChannelStrings(clientIn, c)

	for {
		select {
		case req, ok := <-clientIn:
			{
				if !ok {
					return
				}
				exit, err := handleClientRequest(thisClient, req)
				if exit {
					return
				}
				if err != nil {
					log.Println(err)
					continue
				}
			}
		case <-server.Ctx.Done():
			return // will close the client handle
		}
	}
}

func EntryPoint(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("server needs exactly 1 argument\n" +
			"The port number on which to listen")
	}
	log.Println("Starting server...")
	Run(args[0])
	return nil
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
		Clients: map[string]*clientWrapperTy{},
	}

	stdin := make(chan string, 1)

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

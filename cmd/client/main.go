package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/ZOrfeas/go_chat/cmd/utils"
)

type cliTy struct {
	Id   string
	Conn net.Conn
}

func (cl cliTy) sendBytes(b []byte) error {
	_, err := cl.Conn.Write(append(b, '\n'))
	return err
}
func (cl cliTy) SendString(str string) error {
	return cl.sendBytes([]byte(str))
}

var commands []func(string)

func (cl cliTy) InitHostCommands() {
	exit := func(s string) {
		fmt.Println("disconnect")
		os.Exit(0)
	}
	changeId := func(id string) {
		fmt.Println("change username to " + id)
		cl.Id = id
	}
	commands = []func(string){exit, changeId}
}

func (cl cliTy) ExeHostCommand(idx int, arg string) {
	fmt.Print("Server command: ")
	commands[idx](arg)
}

var client cliTy

func channelStrings(out chan<- string, in io.Reader) {
	reader := bufio.NewReader(in)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		out <- message
	}
}

func initNameFromHost() error {
	if err := client.SendString("/name " + client.Id); err != nil {
		return err
	}
	reader := bufio.NewReader(client.Conn)
	hostResponse, err := reader.ReadString('\n')
	client.Id = utils.RemoveSpace(hostResponse)
	return err
}

func handleUserInput(input string) (err error) {
	input = strings.TrimFunc(input, unicode.IsSpace)
	err = client.SendString(input)
	return
}

func handleHostCommand(cmd string) error {
	fields := strings.Fields(cmd)
	if len(fields) > 2 {
		fmt.Println("Internal error")
		os.Exit(1)
	}
	if len(fields) == 1 {
		fields = append(fields, "")
	}
	idx, err := strconv.Atoi(fields[0])
	if err != nil {
		return err
	}
	client.ExeHostCommand(idx, fields[1])
	return nil
}
func handleHostMessage(msg string) error {
	if strings.HasPrefix(msg, utils.HostCommandIdentifier) {
		msg = strings.TrimPrefix(msg, utils.HostCommandIdentifier)
		return handleHostCommand(msg)
	}
	fmt.Print(msg)
	return nil
}

func Run(connString, id string) {
	id = strings.ReplaceAll(id, " ", "_")
	client.Id = id

	fmt.Println("Dialing host")
	c, err := net.Dial("tcp", connString)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Conn = c
	fmt.Println("Network connection established")

	if err := initNameFromHost(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Host accepted connection")
	fmt.Println("Your username is: " + client.Id)
	fmt.Println("Use '/help' to get a list of available commands")

	client.InitHostCommands()

	stdin := make(chan string, 1)
	host := make(chan string, 1)
	defer close(stdin)
	defer close(host)

	go channelStrings(stdin, os.Stdin)
	go channelStrings(host, client.Conn)

	var exitStatus error
	for {
		select {
		case input := <-stdin:
			exitStatus = handleUserInput(input)
		case message := <-host:
			exitStatus = handleHostMessage(message)
		}
		if exitStatus != nil {
			fmt.Println(exitStatus)
			return
		}
	}

}

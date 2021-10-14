package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode"

	common "github.com/ZOrfeas/go_chat/common/utils"
)

type CliTy struct {
	Id   string
	Conn net.Conn
}

func (cl *CliTy) setName(newName string) {
	if newName == "" {
		fmt.Println("Empty name given")
		return
	}
	cl.Id = newName
}
func (cl *CliTy) sendBytes(b []byte) error {
	_, err := cl.Conn.Write(append(b, '\n'))
	return err
}
func (cl *CliTy) SendString(str string) error {
	return cl.sendBytes([]byte(str))
}
func (cl *CliTy) exeHostCommand(idx common.HostCommand, arg string) {
	fmt.Println("Server command:", idx.String(), "with arg", "'"+arg+"'")
	switch idx {
	case common.Disconnect:
		os.Exit(0)
	case common.ChangeName:
		cl.setName(arg)
	case common.SayName:
		cl.SendString(cl.Id)
	default:
		fmt.Println(idx.String(), " with arg ", arg, "\n", "Not yet implemented")
	}
}

var client *CliTy

const FirstMessagePrefix = "/name "

func initNameFromHost() error {
	if err := client.SendString(FirstMessagePrefix + client.Id); err != nil {
		return err
	}
	reader := bufio.NewReader(client.Conn)
	hostResponse, err := reader.ReadString('\n')
	client.Id = common.RemoveSpace(hostResponse)
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
	client.exeHostCommand(common.HostCommand(idx), fields[1])
	return nil
}
func handleHostMessage(msg string) error {
	if strings.HasPrefix(msg, common.HostCommandIdentifier) {
		msg = strings.TrimPrefix(msg, common.HostCommandIdentifier)
		return handleHostCommand(msg)
	}
	fmt.Print(msg)
	return nil
}

func EntryPoint(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("client needs exactly 2 arguments\n" +
			"The connection string and a username")
	}
	fmt.Println("Starting client...")
	Run(args[0], args[1])
	return nil

}

func Run(connString, id string) {
	client = &CliTy{}
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

	stdin := make(chan string, 1)
	host := make(chan string, 1)

	go common.ChannelStrings(stdin, os.Stdin)
	go common.ChannelStrings(host, client.Conn)

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

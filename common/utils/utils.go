package utils

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type HostCommand int

const (
	Disconnect HostCommand = iota
	ChangeName
	SayName

	HostCommandCount // the count of available host commands
)

// Array from HostCommand enum to string
var hostCommandString = [HostCommandCount]string{
	"Disconnect", "Change-Name", "Say-Name",
}

func (hostCmd HostCommand) String() string {
	if hostCmd >= HostCommandCount || hostCmd < 0 {
		return "Non-existent command ID"
	}
	return hostCommandString[hostCmd]
}

// Map from string to HostCommand enum
var StringHostCommand = func() map[string]HostCommand {
	mapRes := map[string]HostCommand{}
	for i, name := range hostCommandString {
		mapRes[name] = HostCommand(i)
	}
	return mapRes
}

const HostCommandIdentifier = "--COMMAND:"

func ChannelStrings(out chan<- string, in io.Reader) {
	defer close(out)
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

func RemoveSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

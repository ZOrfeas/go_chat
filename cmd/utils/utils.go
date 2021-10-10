package utils

import (
	"unicode"
)

type HostCommand int

const (
	Disconnect HostCommand = iota
	ChangeName
	SayName

	HostCommandCount // the count of available host commands
)

func (hostCmd HostCommand) String() string {
	if hostCmd > HostCommandCount || hostCmd < 0 {
		return "Non-existent command ID"
	}
	return [HostCommandCount]string{
		"Disconnect", "Change-Name", "Say-Name",
	}[hostCmd]
}

const HostCommandIdentifier = "--COMMAND:"

func RemoveSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

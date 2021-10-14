package utils

import (
	"log"
	"strings"

	common "github.com/ZOrfeas/go_chat/common/utils"
)

type service int

const (
	help service = iota
	exit
	name
	list
	msg

	serviceCount
)

// Array from service enum to string
var serviceStrings = [...]string{
	"/help", "/exit", "/name", "/list", "/msg",
}
var providers = [...]func(*clientWrapperTy, string) error{
	helpProvider, exitProvider, nameProvider,
	listProvider, msgProvider,
}

// Map from string to service enum
func (serv service) String() string {
	if serv < 0 || serv >= serviceCount {
		return "INVALID_SERVICE"
	}
	return serviceStrings[serv]
}

var stringServices = func() map[string]service {
	mapRes := map[string]service{}
	for i, name := range serviceStrings {
		mapRes[name] = service(i)
	}
	return mapRes
}()

func exitProvider(thisClient *clientWrapperTy, rest string) error {
	return thisClient.sendCommand(common.Disconnect, "")
}
func nameProvider(thisClient *clientWrapperTy, rest string) error {
	newName := strings.Fields(rest)[0]
	delete(server.Clients, thisClient.Client.Id)
	server.Clients[newName] = thisClient
	return thisClient.sendCommand(common.ChangeName, strings.Fields(rest)[0])
}
func listProvider(thisClient *clientWrapperTy, rest string) error {
	var bd strings.Builder
	bd.WriteString("List of currently active users\n")
	var lineSoFar int
	for key := range server.Clients {
		bd.WriteString(key)
		if lineSoFar == 2 {
			bd.WriteRune('\n')
		} else {
			bd.WriteRune('\t')
		}
		lineSoFar++
	}
	return thisClient.checkAndSend(bd.String())
}
func msgProvider(thisClient *clientWrapperTy, rest string) error {
	words := strings.Fields(rest)
	receiverName := words[0]
	receiver, exists := server.Clients[receiverName]
	if !exists {
		return thisClient.checkAndSend("User " + receiverName + " not found")
	}
	reply := thisClient.Client.Id + " whispers: " +
		strings.TrimPrefix(rest, " "+receiverName+" ")
	return receiver.checkAndSend(reply)
}
func helpProvider(thisClient *clientWrapperTy, rest string) error {
	var bd strings.Builder
	bd.WriteString("List of available services:\n")
	bd.WriteString("\t/help\tgives this message\n")
	bd.WriteString("\t/exit\tdisconnects from server\n")
	bd.WriteString("\t/name\trequests a name change with the first word given\n")
	bd.WriteString("\t/list\treturns a list of currently active usernames\n")
	bd.WriteString("\t/msg\t whispers a message to the username given")
	return thisClient.checkAndSend(bd.String())
}
func broadCastProvider(thisClient *clientWrapperTy, rest string) error {
	log.Println("Client", thisClient.Client.Id, "broadcast", rest)
	for _, clientWrapper := range server.Clients {
		if clientWrapper.Client.Id == thisClient.Client.Id {
			continue
		}
		if err := clientWrapper.checkAndSend(rest); err != nil {
			errMsg := "Failed to send to " + clientWrapper.Client.Id
			log.Println(errMsg)
		}
	}
	return nil
}

func handleClientRequest(thisClient *clientWrapperTy, req string) (bool, error) {
	words := strings.Fields(req)
	serviceRequested, exists := stringServices[words[0]]
	if !exists {
		return false, broadCastProvider(thisClient, req)
	}
	err := providers[int(serviceRequested)](
		thisClient, strings.TrimPrefix(req, words[0]))
	if serviceRequested == exit {
		return true, err
	}
	return false, err

}

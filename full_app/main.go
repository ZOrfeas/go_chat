package main

import (
	"fmt"
	"os"

	client "github.com/ZOrfeas/go_chat/client/utils"
	"github.com/ZOrfeas/go_chat/common/cli"
	server "github.com/ZOrfeas/go_chat/server/utils"
)

func main() {
	parentCli := cli.NewCli(os.Args[0], "Server-Client Cli")

	tokenExists := func(token string, args []string) bool {
		for _, curr := range args {
			if curr == token {
				return true
			}
		}
		return false
	}
	if tokenExists("server", os.Args) && tokenExists("client", os.Args) {
		fmt.Println("client and server commands cannot be used together")
		return
	}

	clientCallback := client.EntryPoint
	serverCallback := server.EntryPoint

	parentCli.AddCommand("client", "chooses client functionality", clientCallback)
	parentCli.AddCommand("server", "chooses server functionality", serverCallback)

	err := parentCli.Run(os.Args[1:], parentCli.Help)
	if err != nil {
		fmt.Println(err)
	}
}

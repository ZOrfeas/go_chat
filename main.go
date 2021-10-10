package main

import (
	"fmt"
	"os"

	client "github.com/ZOrfeas/go_chat/client/utils"
	"github.com/ZOrfeas/go_chat/client_server/cli"
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

	clientCallback := func(args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("client needs exactly 2 arguments\n" +
				"The connection string and a username")
		}
		fmt.Println("Starting client with args: ", args)
		client.Run(args[0], args[1])
		return nil
	}
	serverCallback := func(args []string) error {
		fmt.Println("I am server, these are my args", args)
		return nil
	}

	parentCli.AddCommand("client", "chooses client functionality", clientCallback)
	parentCli.AddCommand("server", "chooses server functionality", serverCallback)

	err := parentCli.Run(os.Args[1:], parentCli.Help)
	if err != nil {
		fmt.Println(err)
	}
}

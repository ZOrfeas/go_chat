package main

import (
	"fmt"
	"os"

	"github.com/ZOrfeas/go_chat/pkg/cli"
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
		fmt.Println("I am client, these are my args", args)
		return nil
	}
	serverCallback := func(args []string) error {
		fmt.Println("I am server, these are my args", args)
		return nil
	}

	parentCli.AddCommand("client", "chooses client functionality", clientCallback)
	parentCli.AddCommand("server", "chooses server functionality", serverCallback)

	parentCli.Run(os.Args[1:])
}

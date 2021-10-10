package main

import (
	"fmt"
	"os"

	"C"

	"github.com/ZOrfeas/go_chat/client/utils"
)

//export RunClient
func RunClient(args []string) {
	if err := utils.EntryPoint(args); err != nil {
		fmt.Println(err)
	}
}

func main() {
	RunClient(os.Args[1:])
}

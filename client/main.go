package main

import (
	"fmt"
	"os"

	"C"

	"github.com/ZOrfeas/go_chat/client/utils"
)

// Can be used externally by dup-ing its stdin and stdout file descriptors
// export RunClient
func RunClient(connString, id string) {
	utils.Run(connString, id)
}

func main() {
	if err := utils.EntryPoint(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

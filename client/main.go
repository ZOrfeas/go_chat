package main

import (
	"fmt"
	"os"

	"github.com/ZOrfeas/go_chat/client/utils"
)

func main() {
	if err := utils.EntryPoint(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

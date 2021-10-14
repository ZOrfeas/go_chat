package main

import "C"
import (
	"fmt"
	"os"
	"strconv"

	"github.com/ZOrfeas/go_chat/server/utils"
)

// Can be used externally by dup-ing its stdin and stdout file
// export RunServer
func RunServer(portNr int) {
	portStr := strconv.FormatInt(int64(portNr), 10)
	utils.Run(portStr)
}

func main() {
	if err := utils.EntryPoint(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

package cli

import (
	"fmt"
	"os"
)

// Showcases a simple usecase of the Cli package
func ExampleCli_Run() {
	testCli := NewCli(os.Args[0], "Just a simple example")
	testCallback := func(args []string) error {
		fmt.Println(args)
		return nil
	}
	testCli.AddCommand("cmd1", "This is cmd1", testCallback)
	testCli.AddCommand("cmd2", "This is cmd2", testCallback)
	testCli.AddCommand("cmd3", "This is cmd3", testCallback)

	testCli.Run([]string{"cmd1", "trsh11", "trsh12", "cmd2", "cmd3", "trsh31"})
	// Output:
	// [trsh11 trsh12]
	// []
	// [trsh31]
}

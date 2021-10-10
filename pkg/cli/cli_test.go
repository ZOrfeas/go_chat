package cli

import (
	"os"
	"strings"
	"testing"
)

// Showcases a simple usecase of the Cli package
func TestCli_Run(t *testing.T) {
	testCli := NewCli(os.Args[0], "Just a simple example")
	bld := strings.Builder{}
	testCallback := func(args []string) error {
		bld.WriteString(strings.Join(args, " "))
		bld.WriteString("-")
		return nil
	}
	testCli.AddCommand("cmd1", "This is cmd1", testCallback)
	testCli.AddCommand("cmd2", "This is cmd2", testCallback)
	testCli.AddCommand("cmd3", "This is cmd3", testCallback)

	testCli.Run([]string{"cmd1", "trsh11", "trsh12", "cmd2", "cmd3", "trsh31"}, testCli.Help)
	if str := bld.String(); str != "trsh11 trsh12--trsh31-" {
		t.Fatal("Did not match: ", str)
	}
}

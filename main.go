package main

import (
	"github.com/rwiteshbera/Hyperion/cli"
	"os"
)

func main() {

	defer os.Exit(0)
	cmd := cli.CLI{}
	cmd.RunCLI()

}

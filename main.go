package main

import (
	"os"

	"github.com/rwiteshbera/Hyperion/cli"
)

func main() {

	defer os.Exit(0)
	cmd := cli.CLI{}
	cmd.RunCLI()
}

// Get IST

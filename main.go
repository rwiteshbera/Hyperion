package main

import (
	"os"

	"github.com/rwiteshbera/Blockchain-Go/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CLI{}
	cmd.RunCLI()
}

// Get IST

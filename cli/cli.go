package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/rwiteshbera/Blockchain-Go/wallet"
)

type CLI struct{}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}
}

func (cli *CLI) create() {
	wallets, _ := wallet.InitWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()
	fmt.Printf("%s\n", address)
}

func (cli *CLI) list() {
	wallets, _ := wallet.InitWallets()
	addresses := wallets.GetAllAddresses()
	// wallets.SaveFile()
	for _, address := range addresses {
		fmt.Printf("%s\n", address)
	}
}

func (cli *CLI) RunCLI() {
	cli.validateArgs()

	CreateWalletCMD := flag.NewFlagSet("--create", flag.ExitOnError)
	ListAddressesCMD := flag.NewFlagSet("--list", flag.ExitOnError)

	switch os.Args[1] {
	case "--create":
		err := CreateWalletCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "--list":
		err := ListAddressesCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		runtime.Goexit()
	}

	if CreateWalletCMD.Parsed() {
		cli.create()
	}
	if ListAddressesCMD.Parsed() {
		cli.list()
	}

}

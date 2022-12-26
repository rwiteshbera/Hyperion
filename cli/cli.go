package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/rwiteshbera/Hyperion/wallet"
)

type CLI struct{}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}
}

func (cli *CLI) create() {
	wallets, _ := wallet.InitWallets()      // Init saved wallets from file
	newWallet := wallet.CreateWallet(0)     // Create new wallet
	address := wallets.AddWallet(newWallet) // add it wallets file
	wallets.SaveFile()                      // save it
	fmt.Printf("%s\n", address)
}

func (cli *CLI) list() {
	wallets, _ := wallet.InitWallets()
	addresses := wallets.GetAllAddresses()
	// wallets.SaveFile()
	for _, address := range addresses {
		fmt.Printf("%s : %f\n", address, wallets.GetWallet(address).Balance)
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

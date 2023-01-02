package cli

import (
	"flag"
	"fmt"
	"github.com/rwiteshbera/Hyperion/network"
	"github.com/rwiteshbera/Hyperion/wallet"
	"log"
	"os"
	"runtime"
)

type CLI struct{}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}
}

// Create a new wallet
func (cli *CLI) createWallet() {
	wallets, _ := wallet.LoadWallets() // Init saved wallets from file
	newWallet := wallet.CreateWallet() // Create new wallet
	wallets.AddWallet(newWallet)       // add it to wallets
	wallets.SaveFile()                 // save it

	privatekey, publickey, err := newWallet.GetKeyValuePair()
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Printf("Wallet Address : %s\nPrivate Key : %s\nPublic Key : %s\n", newWallet.GetWalletAddress(), privatekey, publickey)
}

// List all the save wallets with balances
func (cli *CLI) listWallets() {
	wallets, _ := wallet.LoadWallets()
	addresses := wallets.GetAllAddresses()
	// wallets.SaveFile()
	for _, address := range addresses {
		fmt.Printf("%s\n", address)
	}
}

// Start the blockchain server on specific port
func (cli *CLI) startBlockchainServer(port *string) {
	network.Server(port)
}

func (cli *CLI) RunCLI() {
	cli.validateArgs()

	CreateWalletCMD := flag.NewFlagSet("--create", flag.ExitOnError) // --create
	ListAddressesCMD := flag.NewFlagSet("--list", flag.ExitOnError)  // --list
	ServerStartupPort := flag.NewFlagSet("--port", flag.ExitOnError) // --port <MENTION_PORT_NUMBER>

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
	case "--port":
		err := ServerStartupPort.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		runtime.Goexit()
	}

	if CreateWalletCMD.Parsed() {
		cli.createWallet()
	}
	if ListAddressesCMD.Parsed() {
		cli.listWallets()
	}
	if ServerStartupPort.Parsed() {
		port := ServerStartupPort.String("port", "", "The port to listen on")
		err := ServerStartupPort.Parse(os.Args[1:])
		if err != nil {
			return
		}
		cli.startBlockchainServer(port)
	}

}

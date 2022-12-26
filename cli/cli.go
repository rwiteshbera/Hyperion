package cli

import (
	"flag"
	"fmt"
	"github.com/rwiteshbera/Hyperion/network"
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

// Create a new wallet
func (cli *CLI) createWallet() {
	wallets, _ := wallet.InitWallets()      // Init saved wallets from file
	newWallet := wallet.CreateWallet(0)     // Create new wallet
	address := wallets.AddWallet(newWallet) // add it to wallets
	wallets.SaveFile()                      // save it
	fmt.Printf("%s\n", address)
}

// List all the save wallets with balances
func (cli *CLI) listWallets() {
	wallets, _ := wallet.InitWallets()
	addresses := wallets.GetAllAddresses()
	// wallets.SaveFile()
	for _, address := range addresses {
		fmt.Printf("%s : %f\n", address, wallets.GetWallet(address).Balance)
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
		ServerStartupPort.Parse(os.Args[1:])
		cli.startBlockchainServer(port)
	}

}

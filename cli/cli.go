package cli

import (
	"flag"
	"fmt"
	"github.com/rwiteshbera/Hyperion/network"
	"github.com/rwiteshbera/Hyperion/wallet"
	"log"
	"os"
	"runtime"
	"strings"
)

type CLI struct{}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		runtime.Goexit()
	}
}

// Create a new wallet
func (cli *CLI) createWallet() {
	newWallet := wallet.CreateWallet() // Create new wallet

	privatekey, publickey, err := newWallet.GetKeyValuePair() // Get the privatekey and publick key in string format
	if err != nil {
		log.Panic(err.Error())
	}
	database := wallet.ConnectBitcask()
	defer database.Close()

	// Wallet data will be stored in key - value pair format
	// Key: Wallet Address, Value: Private Key and Public Key seperated by ' ' single space (privateKey+" "+publicKey)
	err = database.Put([]byte(newWallet.GetWalletAddress()), []byte(privatekey+" "+publickey))
	if err != nil {
		log.Panic(err.Error())
	}

	// Print new wallet data in console
	fmt.Printf("Wallet Address : %s\nPrivate Key : %s\nPublic Key : %s\n", newWallet.GetWalletAddress(), privatekey, publickey)
	fmt.Printf("\nWarning: Never disclose the private key. Anyone with your private key can steal your assets.\nIt is important to ensure that you have securely saved your private key and public key.\nThe \"--list\" command will display only the wallet addresses.")
}

// List all the save wallets with balances
func (cli *CLI) listWallets() {
	database := wallet.ConnectBitcask()
	defer database.Close()
	walletAddresses := database.Keys()
	for i := range walletAddresses {
		fmt.Printf("%s\n", i)
	}
}

// Get Private Key
func (cli *CLI) getPrivatePublicKey(walletAddress *string) {
	database := wallet.ConnectBitcask()
	defer database.Close()

	value, err := database.Get([]byte(*walletAddress))
	if err != nil {
		log.Panic(err.Error())
	}

	// Private Key and Public Key are stored in a single string seperated by space. This function helps to split it.
	keyValue := strings.Split(string(value), " ")

	// Returning Private Key + Public Key
	fmt.Println("Private Key: ", keyValue[0])
	fmt.Println("Public Key: ", keyValue[1])
	fmt.Println()
	fmt.Println("Warning: Never disclose the private key. Anyone with your private key can steal your assets.")
}

// Start the blockchain server on specific port
func (cli *CLI) startBlockchainServer(port *string) {
	network.Server(port)
}

func (cli *CLI) RunCLI() {
	cli.validateArgs()

	CreateWalletCMD := flag.NewFlagSet("--create", flag.ExitOnError)     // --create
	ListAddressesCMD := flag.NewFlagSet("--list", flag.ExitOnError)      // --list
	GetPrivatePublicKeyCMD := flag.NewFlagSet("--get", flag.ExitOnError) // --get <WALLET_ADDRESS>
	ServerStartupPort := flag.NewFlagSet("--port", flag.ExitOnError)     // --port <MENTION_PORT_NUMBER>

	switch os.Args[1] {
	case "--create":
		err := CreateWalletCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	case "--list":
		err := ListAddressesCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	case "--get":
		err := GetPrivatePublicKeyCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
		}
	case "--port":
		err := ServerStartupPort.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err.Error())
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
	if GetPrivatePublicKeyCMD.Parsed() {
		walletAddress := GetPrivatePublicKeyCMD.String("get", "", "Wallet Address")
		err := GetPrivatePublicKeyCMD.Parse(os.Args[1:])
		if err != nil {
			log.Panic(err.Error())
		}
		cli.getPrivatePublicKey(walletAddress)
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

package main

import (
	"fmt"

	"github.com/rwiteshbera/Blockchain-Go/blockchain"
	"github.com/rwiteshbera/Blockchain-Go/wallet"
)

func main() {
	// defer os.Exit(0)
	// cmd := cli.CLI{}
	// cmd.RunCLI()

	wallet1 := wallet.CreateWallet(0)  // sender
	wallet2 := wallet.CreateWallet(50) // recipient
	private1, public1 := wallet1.GetKeyValuePair()
	_, publicAddress := wallet1.GetKeyValuePairInString()
	_, recipientAddress := wallet2.GetKeyValuePairInString()

	fmt.Println(wallet1.GetBalance())
	fmt.Println(wallet2.GetBalance())

	tr1 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	sig1 := tr1.GenerateSignature()
	tr1.Mining(public1, sig1, wallet1, wallet2)

	fmt.Println(wallet1.GetBalance())
	fmt.Println(wallet2.GetBalance())

}

// Get IST

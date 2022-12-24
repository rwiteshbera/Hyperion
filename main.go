package main

import (
	"os"

	"github.com/rwiteshbera/Blockchain-Go/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CLI{}
	cmd.RunCLI()

	// chain := blockchain.InitBlockchain()

	// wallet1 := wallet.CreateWallet(30) // sender
	// wallet2 := wallet.CreateWallet(50) // recipient
	// walletM := wallet.CreateWallet(0)  // miner

	// fmt.Println(walletM.Balance)

	// private1, public1 := wallet1.GetKeyValuePair()
	// _, publicAddress := wallet1.GetKeyValuePairInString()
	// _, recipientAddress := wallet2.GetKeyValuePairInString()

	// tr1 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	// sig1 := tr1.GenerateSignature()
	// chain.Mining(tr1, public1, sig1, wallet1, wallet2, walletM)

	// fmt.Println(walletM.Balance)

	// tr2 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	// sig2 := tr2.GenerateSignature()
	// chain.Mining(tr2, public1, sig2, wallet1, wallet2, walletM)

	// fmt.Println(walletM.Balance)

	// tr3 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	// sig3 := tr3.GenerateSignature()
	// chain.Mining(tr3, public1, sig3, wallet1, wallet2, walletM)

	// fmt.Println(walletM.Balance)

	// tr4 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	// sig4 := tr4.GenerateSignature()
	// chain.Mining(tr4, public1, sig4, wallet1, wallet2, walletM)

	// fmt.Println(walletM.Balance)

	// tr5 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	// sig5 := tr5.GenerateSignature()
	// chain.Mining(tr4, public1, sig5, wallet1, wallet2, walletM)

	// fmt.Println(walletM.Balance)

	// tr6 := blockchain.NewTransaction(private1, public1, publicAddress, recipientAddress, 5)
	// sig6 := tr6.GenerateSignature()
	// chain.Mining(tr4, public1, sig6, wallet1, wallet2, walletM)

	// fmt.Println(walletM.Balance)

	// chain.ShowBlockchain()

}

// Get IST

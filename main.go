package main

import (
	"fmt"
	"github.com/rwiteshbera/Hyperion/blockchain"
	"github.com/rwiteshbera/Hyperion/wallet"
)

func main() {

	//defer os.Exit(0)
	//cmd := cli.CLI{}
	//cmd.RunCLI()

	bcs := blockchain.InitBlockchain()
	wallet1 := wallet.CreateWallet(10)
	wallet2 := wallet.CreateWallet(0)

	walletM := wallet.CreateWallet(0)

	private1, public1 := wallet1.GetKeyValuePair()
	_, sender := wallet1.GetKeyValuePairInString()
	_, recipient := wallet2.GetKeyValuePairInString()

	fmt.Println(wallet1.GetBalance())
	fmt.Println(wallet2.GetBalance())

	bcs.NewTransaction(private1, public1, sender, recipient, 1)
	bcs.NewTransaction(private1, public1, sender, recipient, 2)
	bcs.Mining(wallet1, wallet2, walletM)
	bcs.Mining(wallet1, wallet2, walletM)

	fmt.Println(wallet1.GetBalance())
	fmt.Println(wallet2.GetBalance())

}

// Get IST

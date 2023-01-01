package main

import (
	"fmt"
	"github.com/rwiteshbera/Hyperion/blockchain"
	"github.com/rwiteshbera/Hyperion/wallet"
)

func test() {
	bcs := blockchain.InitBlockchain()

	wallet1 := wallet.CreateWallet()
	wallet2 := wallet.CreateWallet()

	private1, public1 := wallet1.GetKeyValuePair()
	sender := wallet1.GetWalletAddress()
	recipient := wallet2.GetWalletAddress()

	bcs.NewTransaction(private1, public1, sender, recipient, 1)
	bcs.NewTransaction(private1, public1, sender, recipient, 2)

	bcs.StartMining()
	bcs.StartMining()

	fmt.Println(wallet1.GetWalletBalance(bcs))
	fmt.Println(wallet2.GetWalletBalance(bcs))
}

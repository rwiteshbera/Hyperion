package main

import (
	"github.com/rwiteshbera/Hyperion/blockchain"
	"github.com/rwiteshbera/Hyperion/wallet"
)

func main() {

	//defer os.Exit(0)
	//cmd := cli.CLI{}
	//cmd.RunCLI()

	bcs := blockchain.InitBlockchain()
	wallet1 := wallet.CreateWallet(3)
	wallet2 := wallet.CreateWallet(0)

	private1, public1 := wallet1.GetKeyValuePair()
	sender := wallet1.GetWalletAddress()
	recipient := wallet2.GetWalletAddress()

	bcs.NewTransaction(private1, public1, sender, recipient, 1)
	bcs.Mining(wallet1, wallet2)

	bcs.NewTransaction(private1, public1, sender, recipient, 2)
	bcs.Mining(wallet1, wallet2)

	bcs.NewTransaction(private1, public1, sender, recipient, 1)
	bcs.Mining(wallet1, wallet2)

	bcs.NewTransaction(private1, public1, sender, recipient, 2)
	bcs.Mining(wallet1, wallet2)

	bcs.ShowBlockchain()
}

// Get IST

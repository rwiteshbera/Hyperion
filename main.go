package main

import (
	"github.com/rwiteshbera/Hyperion/cli"
	"os"
)

func main() {

	defer os.Exit(0)
	cmd := cli.CLI{}
	cmd.RunCLI()
	//bcs := blockchain.InitBlockchain()
	//
	//wallet1 := wallet.CreateWallet()
	//wallet2 := wallet.CreateWallet()
	//
	//private1, public1, _ := wallet1.GetKeyValuePair()
	//sender := wallet1.GetWalletAddress()
	//recipient := wallet2.GetWalletAddress()
	//
	//bcs.NewTransaction(private1, public1, sender, recipient, 1)
	//bcs.NewTransaction(private1, public1, sender, recipient, 2)
	//
	//bcs.StartMining()
	//bcs.StartMining()
	//
	//fmt.Println(wallet1.GetWalletBalance(bcs))
	//fmt.Println(wallet2.GetWalletBalance(bcs))
}

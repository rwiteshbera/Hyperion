package main

import (
	"fmt"
	"strconv"

	"github.com/rwiteshbera/Blockchain-Go/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()
	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	for _, block := range chain.Blocks {
		fmt.Printf("Block Number : %x\n", block.BlockNumber)
		fmt.Printf("Nonce : %x\n", block.Nonce)
		fmt.Printf("Data : %x\n", block.Data)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Previous Hash : %x\n", block.PreviousHash)
		fmt.Printf("TimeStamp : %x\n", block.TimeStamp)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW : %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

// Get IST

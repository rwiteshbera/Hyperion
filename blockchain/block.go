package blockchain

import (
	"fmt"
	"time"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	BlockNumber         int
	Nonce               int
	TransactionsInBlock []*Transaction
	Hash                []byte
	PreviousHash        []byte
	TimeStamp           int64
}

// Create New Block
func createBlock(NewBlockNumber int, transactionsData []*Transaction, previousHash []byte) *Block {
	block := &Block{
		BlockNumber:         NewBlockNumber,
		Nonce:               0,
		TransactionsInBlock: transactionsData,
		Hash:                []byte{},
		PreviousHash:        previousHash,
		TimeStamp:           time.Now().Unix(),
	}
	proofOfWork := newProof(block)

	nonce, hash := proofOfWork.run()

	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

// Add New Block to the chain
func (blockchain *Blockchain) AddBlock(transactionsData []*Transaction, minerAddress string) {
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	NewBlockNumber := len(blockchain.Blocks) + 1
	NewBlock := createBlock(NewBlockNumber, transactionsData, previousBlock.Hash)
	blockchain.Blocks = append(blockchain.Blocks, NewBlock)
}

// Build The First Block // Genesis Block
// https://en.bitcoin.it/wiki/Genesis_block
func Genesis() *Block {
	return createBlock(1, []*Transaction{}, []byte{})
}

/*
	The InitBlockchain function is defining a new function called InitBlockchain that returns a pointer to a Blockchain struct. The Blockchain struct contains a slice of pointers to Block structs, and the slice is initialized with a single Block struct that is returned by a function called Genesis.

Genesis function creates and returns the first block in the blockchain, often referred to as the "genesis block." The purpose of the InitBlockchain function is to create a new instance of a blockchain data structure and initialize it with the genesis block.
*/

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func (blockchain *Blockchain) ShowBlockchain() {
	for _, e := range blockchain.Blocks {
		fmt.Println(e)
		fmt.Println()
	}
}

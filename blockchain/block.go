package blockchain

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Blockchain struct {
	Mempool           []*Transaction // List of confirmed transactions to be added on next block
	TransactionsQueue []*Transaction // List of transactions to be verified
	Blocks            []*Block
	mux               sync.Mutex
}

type Block struct {
	BlockNumber         int
	Nonce               int
	TransactionsInBlock []*Transaction
	Hash                []byte
	PreviousHash        []byte
	TimeStamp           string
}

// Create New Block
func createBlock(NewBlockNumber int, transactionsData []*Transaction, previousHash []byte) *Block {
	block := &Block{
		BlockNumber:         NewBlockNumber,
		Nonce:               0,
		TransactionsInBlock: transactionsData,
		Hash:                []byte{},
		PreviousHash:        previousHash,
		TimeStamp:           time.Now().Format(time.RFC3339),
	}

	proofOfWork := newProof(block)
	nonce, hash := proofOfWork.run()
	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

// Add New Block to the chain
func (chain *Blockchain) AddBlock(transactionsData []*Transaction) *Block {
	previousBlock := chain.Blocks[len(chain.Blocks)-1]
	NewBlockNumber := len(chain.Blocks) + 1
	NewBlock := createBlock(NewBlockNumber, transactionsData, previousBlock.Hash)
	chain.Blocks = append(chain.Blocks, NewBlock)
	return NewBlock
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
	return &Blockchain{Blocks: []*Block{Genesis()}}
}

func (chain *Blockchain) ShowBlockchain() {
	fmt.Printf("Genesis Block : ")
	for _, e := range chain.Blocks {
		blockTime := e.TimeStamp
		t, err := time.Parse(time.RFC3339, blockTime)
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("%d => %d %s %d : %d:%d:%d\n", e.BlockNumber, t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
		fmt.Println()
	}
}

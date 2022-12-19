package blockchain

import (
	"time"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	BlockNumber  int
	Nonce        int
	Data         []byte
	Hash         []byte
	PreviousHash []byte
	TimeStamp    int64
}

// Create New Block
func CreateBlock(NewBlockNumber int, data string, previousHash []byte) *Block {
	block := &Block{
		BlockNumber:  NewBlockNumber,
		Nonce:        0,
		Data:         []byte(data),
		Hash:         []byte{},
		PreviousHash: previousHash,
		TimeStamp:    time.Now().Unix(),
	}
	proofOfWork := NewProof(block)

	nonce, hash := proofOfWork.Run()

	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

// Add New Block to the chain
func (blockchain *Blockchain) AddBlock(data string) {
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	NewBlockNumber := len(blockchain.Blocks) + 1
	NewBlock := CreateBlock(NewBlockNumber, data, previousBlock.Hash)
	blockchain.Blocks = append(blockchain.Blocks, NewBlock)
}

// Build The First Block // Genesis Block
// https://en.bitcoin.it/wiki/Genesis_block
func Genesis() *Block {
	return CreateBlock(1, "Genesis", []byte{})
}

/*
	The InitBlockchain function is defining a new function called InitBlockchain that returns a pointer to a Blockchain struct. The Blockchain struct contains a slice of pointers to Block structs, and the slice is initialized with a single Block struct that is returned by a function called Genesis.

Genesis function creates and returns the first block in the blockchain, often referred to as the "genesis block." The purpose of the InitBlockchain function is to create a new instance of a blockchain data structure and initialize it with the genesis block.
*/

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

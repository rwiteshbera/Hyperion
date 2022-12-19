package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	BlockNumber  []byte
	Data         []byte
	Hash         []byte
	PreviousHash []byte
	timestamp    int64
}

// Generate New Hash
func (block *Block) GenerateHash() {
	info := bytes.Join([][]byte{block.Data, block.PreviousHash}, []byte{})
	hash := sha256.Sum256(info)
	block.Hash = hash[:]
}

// Create New Block
func CreateBlock(data string, previousHash []byte) *Block {
	block := &Block{
		BlockNumber:  []byte{},
		Data:         []byte(data),
		Hash:         []byte{},
		PreviousHash: previousHash,
		timestamp:    time.Now().Unix(),
	}
	block.GenerateHash()
	return block
}

// Add New Block to the chain
func (blockchain *Blockchain) AddBlock(data string) {
	previousBlock := blockchain.blocks[len(blockchain.blocks)-1]
	NewBlock := CreateBlock(data, previousBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, NewBlock)
}

// Build The First Block // Genesis Block
// https://en.bitcoin.it/wiki/Genesis_block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

/*
	The InitBlockchain function is defining a new function called InitBlockchain that returns a pointer to a Blockchain struct. The Blockchain struct contains a slice of pointers to Block structs, and the slice is initialized with a single Block struct that is returned by a function called Genesis.

Genesis function creates and returns the first block in the blockchain, often referred to as the "genesis block." The purpose of the InitBlockchain function is to create a new instance of a blockchain data structure and initialize it with the genesis block.
*/

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func main() {
	Blockchain := InitBlockchain()
	Blockchain.AddBlock("First Block after Genesis")
	Blockchain.AddBlock("Second Block after Genesis")
	Blockchain.AddBlock("Third Block after Genesis")

	for _, block := range Blockchain.blocks {
		fmt.Printf("Block Number : %x\n", block.BlockNumber)
		fmt.Printf("Data : %x\n", block.Data)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Previous Hash : %x\n", block.PreviousHash)
		fmt.Printf("TimeStamp : %x\n", block.timestamp)
		fmt.Println()
	}
}

// Temporary Function
func ConvertToIST(value int64) time.Time {
	timestamp := int64(value)

	// Load the IST location
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Unix(timestamp, 0)
	}

	// Convert the timestamp to a time.Time value in the IST location
	return time.Unix(timestamp, 0).In(loc)
}

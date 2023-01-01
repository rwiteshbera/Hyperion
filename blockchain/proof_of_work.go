package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

// Take the block data

// Create a counter / nonce which starts at 0

// Create a hash of the data + nonce

// Check if the hash is less than or equal to target hash / check valid hash

const MiningDifficulty = 10

type ProofOfWork struct {
	Block      *Block
	TargetHash *big.Int
}

func newProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-MiningDifficulty))

	pow := &ProofOfWork{Block: block, TargetHash: target}
	return pow
}

func (pow *ProofOfWork) initData(nonce int) []byte {
	// Serialize each transaction in the slice into a slice of bytes
	transactionBytes := make([][]byte, len(pow.Block.TransactionsInBlock))
	for i, transaction := range pow.Block.TransactionsInBlock {
		transactionBytes[i] = transaction.Serialize()
	}

	// Concatenate all the transaction bytes into a single slice
	allTransactionBytes := bytes.Join(transactionBytes, []byte{})

	data := bytes.Join(
		[][]byte{
			pow.Block.PreviousHash,
			allTransactionBytes,
			toHex(int64(nonce)),
			toHex(int64(MiningDifficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) run() (int, []byte) {
	var intHash big.Int
	var hash []byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.initData(nonce)
		hashInLoop := sha256.Sum256(data)
		intHash.SetBytes(hashInLoop[:])

		if intHash.Cmp(pow.TargetHash) == -1 {
			// fmt.Printf("\r%x", hash)
			hash = hashInLoop[:]
			break
		} else {
			nonce++
		}
	}
	return nonce, hash
}

/*
This function converts an int64 (a 64-bit signed integer) to a slice of bytes in big-endian order.

bytes := ToHex(1234567890)
fmt.Println(bytes) // Output: [72, 144, 173, 205, 1, 0, 0, 0]
*/
func toHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

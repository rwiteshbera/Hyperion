package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

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

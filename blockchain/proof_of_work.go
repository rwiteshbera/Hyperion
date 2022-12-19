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

const Difficulty = 8

type ProofOfWork struct {
	Block      *Block
	TargetHash *big.Int
}

func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{Block: block, TargetHash: target}
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PreviousHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash []byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
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

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.TargetHash) == -1
}

/*
This function converts an int64 (a 64-bit signed integer) to a slice of bytes in big-endian order. The function first creates a new buffer using the bytes.Buffer type from the bytes package. It then writes the value of num to the buffer using the binary.Write function from the binary package. The binary.BigEndian constant specifies that the value should be written in big-endian order, which means that the most significant byte (the "big end") comes first, followed by the second most significant byte, and so on.

If the call to binary.Write succeeds, the function returns the slice of bytes contained in the buffer. If an error occurs, the function will panic with the error message.

bytes := ToHex(1234567890)
fmt.Println(bytes) // Output: [72, 144, 173, 205, 1, 0, 0, 0]
*/
func ToHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

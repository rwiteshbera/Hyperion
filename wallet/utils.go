package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
)

const (
	checkSumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	privateKey    *ecdsa.PrivateKey
	publicKey     *ecdsa.PublicKey
	walletAddress string
}

/*
	Generate ECDSA New Private Key (32bytes) - Public-Key (64bytes) Pair

Return : privateKey(ecdsa.PrivateKey), PublicKey(ecdsa.PublicKey)
*/
func generateNewKeyPair() (ecdsa.PrivateKey, ecdsa.PublicKey) {
	curve := elliptic.P256()

	private, _ := ecdsa.GenerateKey(curve, rand.Reader)

	public := private.PublicKey

	return *private, public
}

// Perform SHA-256 twice on Public Key
func publicKeyHash(publicKeyX []byte, publicKeyY []byte) []byte {
	CombinedPublicHash := append(publicKeyX, publicKeyY...)
	pubHash := sha256.Sum256(CombinedPublicHash)
	pubHash2 := sha256.Sum256(pubHash[:])
	return pubHash2[:]
}

// Perform SHA-256 twice on Versioned-Hash and take only first 4 bytes
func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	return second[:checkSumLength]
}

// Encode byte string to base58
func base58Encode(publicKeyHash []byte) []byte {
	base58PubKey := base58.Encode(publicKeyHash)

	return []byte(base58PubKey)
}

// Decode byte string to base58
// func base58Decode(input []byte) []byte {
// 	decode := base58.Decode(string(input))
// 	return decode
// }

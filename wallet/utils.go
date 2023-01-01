package wallet

import (
	"crypto/aes"
	"crypto/cipher"
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
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
	WalletAddress string
}

/*
	Generate ECDSA New Private Key (32bytes) - Public-Key (64bytes) Pair

Return : PrivateKey(ecdsa.PrivateKey), PublicKey(ecdsa.PublicKey)
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

func encryptWalletData(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate a random initialization vector.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Encrypt the plaintext using AES in CTR mode.
	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// Concatenate the initialization vector and ciphertext into a single byte slice.
	encrypted := append(iv, ciphertext...)
	return encrypted, nil
}

func decryptWalletData(key []byte, encrypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Split the encrypted byte slice into the initialization vector and ciphertext.
	iv, ciphertext := encrypted[:aes.BlockSize], encrypted[aes.BlockSize:]

	// Decrypt the ciphertext using AES in CTR mode.
	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

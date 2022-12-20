package transactions

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
)

type Transaction struct {
	senderPrivateKey       *ecdsa.PrivateKey
	senderPublicKey        *ecdsa.PublicKey
	senderWalletAddress    string
	recipientWalletAddress string
	value                  float32
}

type Signature struct {
	R *big.Int
	S *big.Int
}

// This function creates a new transaction object with the given parameters
func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{senderPrivateKey: privateKey, senderPublicKey: publicKey, senderWalletAddress: sender, recipientWalletAddress: recipient, value: value}
}

// Generating a signature for the transaction.
func (transaction *Transaction) GenerateSignature() *Signature {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.senderPrivateKey, h[:])

	return &Signature{R: r, S: s}
}

// Converting the transaction object to a JSON object.
func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.senderWalletAddress,
		Recipient: transaction.recipientWalletAddress,
		Value:     transaction.value,
	})
}

// Converting the signature to a string.
func (signature *Signature) Signature() string {
	return fmt.Sprintf("%x%x", signature.R.Bytes(), signature.S.Bytes())
}

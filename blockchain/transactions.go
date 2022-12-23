package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/rwiteshbera/Blockchain-Go/wallet"
)

type Transaction struct {
	SenderPrivateKey       *ecdsa.PrivateKey
	SenderPublicKey        *ecdsa.PublicKey
	SenderWalletAddress    string
	RecipientWalletAddress string
	Value                  float32
}

type Signature struct {
	R *big.Int
	S *big.Int
}

const (
	MINING_SNEDER = "THE BLOCKCHAIN"
	MINING_REWARD = 1.0
	MEMPOOL_SIZE  = 2
)

type Mempool struct {
	UnconfirmedTransactions []*Transaction
}

// This function creates a new transaction object with the given parameters
func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, Value float32) *Transaction {
	return &Transaction{SenderPrivateKey: privateKey, SenderPublicKey: publicKey, SenderWalletAddress: sender, RecipientWalletAddress: recipient, Value: Value}
}

// Transfer Balance
func (transaction *Transaction) Transfer(sender *wallet.Wallet, recipient *wallet.Wallet) bool {
	if sender.Balance == 0 {
		fmt.Println("No coin")
		return false
	} else if sender.Balance >= transaction.Value {
		recipient.Balance += transaction.Value
		sender.Balance -= transaction.Value
		fmt.Println("Transferred")
		return true
	} else if sender.Balance < transaction.Value {
		fmt.Println("Not enough balance")
		return false
	}
	return false
}

// Generating a signature for the transaction.
func (transaction *Transaction) GenerateSignature() *Signature {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.SenderPrivateKey, h[:])

	return &Signature{R: r, S: s}
}

// ValidateSignature takes a public key, a transactionData, and a signature, and returns true if the
// signature is valid for the message and public key
func ValidateSignature(signerPublicKey *ecdsa.PublicKey, transaction *Transaction, signature *Signature) bool {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(signerPublicKey, h[:], signature.R, signature.S)
}

// Mining
func (transaction *Transaction) Mining(signerPublicKey *ecdsa.PublicKey, signature *Signature, sender *wallet.Wallet, recipient *wallet.Wallet) {
	valid := ValidateSignature(signerPublicKey, transaction, signature)

	// If the transaction is valid, add it to mempool
	if valid {
		transaction.Transfer(sender, recipient)
		mempool := Mempool{}
		mempool.UnconfirmedTransactions = append(mempool.UnconfirmedTransactions, transaction)

	} else {
		log.Panic("Invalid Transaction")
	}
}

// Converting the transaction object to a JSON object.
func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Value     float32 `json:"Value"`
	}{
		Sender:    transaction.SenderWalletAddress,
		Recipient: transaction.RecipientWalletAddress,
		Value:     transaction.Value,
	})
}

// Converting the signature to a string.
func (signature *Signature) Signature() string {
	return fmt.Sprintf("%x%x", signature.R.Bytes(), signature.S.Bytes())
}

// The Serialize method is a way to convert the data stored in a Transaction object into a slice of bytes that can be stored or transmitted. The serialized data contains all of the information needed to recreate the original Transaction object, including the values of its fields.
func (t *Transaction) Serialize() []byte {

	var buffer bytes.Buffer

	privKeyBytes, _ := x509.MarshalECPrivateKey(t.SenderPrivateKey)
	binary.Write(&buffer, binary.BigEndian, privKeyBytes)

	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(t.SenderPublicKey)
	binary.Write(&buffer, binary.BigEndian, pubKeyBytes)

	binary.Write(&buffer, binary.BigEndian, []byte(t.SenderWalletAddress))

	binary.Write(&buffer, binary.BigEndian, []byte(t.RecipientWalletAddress))

	binary.Write(&buffer, binary.BigEndian, t.Value)

	return buffer.Bytes()
}

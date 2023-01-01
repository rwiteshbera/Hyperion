package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"
)

const (
	TransactionsToStoreInBlock = 2 // How many transactions you want to store in a single block
)

type Transaction struct {
	SenderPrivateKey       *ecdsa.PrivateKey
	SenderPublicKey        *ecdsa.PublicKey
	SenderWalletAddress    string
	RecipientWalletAddress string
	Value                  float32
	Signature              *Signature
}

type Signature struct {
	R *big.Int
	S *big.Int
}

// This function creates a new transaction object with the given parameters
func (chain *Blockchain) NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, Value float32) *Transaction {
	t := &Transaction{SenderPrivateKey: privateKey, SenderPublicKey: publicKey, SenderWalletAddress: sender, RecipientWalletAddress: recipient, Value: Value}
	t.Signature = GenerateSignature(t)
	chain.TransactionsQueue = append(chain.TransactionsQueue, t)
	return t
}

// Generating a signature for the transaction.
func GenerateSignature(transaction *Transaction) *Signature {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256(m)
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.SenderPrivateKey, h[:])

	return &Signature{R: r, S: s}
}

// ValidateSignature takes a public key, a transactionData, and a signature, and returns true if the
// signature is valid for the message and public key
func ValidateSignature(signerPublicKey *ecdsa.PublicKey, transaction *Transaction, signature *Signature) bool {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256(m)
	return ecdsa.Verify(signerPublicKey, h[:], signature.R, signature.S)
}

// Mining
func (chain *Blockchain) mine() {
	chain.mux.Lock()
	defer chain.mux.Unlock()

	// Check if the transactions queue is empty or not?
	if len(chain.TransactionsQueue) == 0 {
		fmt.Println("No transaction in queue")
		return
	}

	// Take first transaction for validation and also remove it from queue
	var tr1 *Transaction
	var err error
	tr1, chain.TransactionsQueue, err = dequeue(chain.TransactionsQueue)
	if err != nil {
		log.Panic(err)
	}

	transactionWithoutSignature := &Transaction{
		SenderPrivateKey:       tr1.SenderPrivateKey,
		SenderPublicKey:        tr1.SenderPublicKey,
		SenderWalletAddress:    tr1.SenderWalletAddress,
		RecipientWalletAddress: tr1.RecipientWalletAddress,
		Value:                  tr1.Value,
	}
	valid := ValidateSignature(transactionWithoutSignature.SenderPublicKey, transactionWithoutSignature, tr1.Signature)

	// If the transaction is valid, add it to Mempool
	if valid {
		chain.Mempool = enqueue(chain.Mempool, tr1)

		if len(chain.Mempool) == TransactionsToStoreInBlock {
			// Adding transactions to block
			chain.AddBlock(chain.Mempool)

			// Removing the first two transactions from the mem-pool as it is already added to block
			chain.Mempool = chain.Mempool[TransactionsToStoreInBlock:]

			return
		}
	} else {
		fmt.Println("Transaction Failed")
	}

}

// Start the mining process
func (chain *Blockchain) StartMining() {
	chain.mine()
	_ = time.AfterFunc(time.Second*5, chain.StartMining)
}

// Converting the signature to a string.
func (signature *Signature) Signature() string {
	return fmt.Sprintf("%x%x", signature.R.Bytes(), signature.S.Bytes())
}

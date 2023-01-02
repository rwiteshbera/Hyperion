package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
	HashId                 string
	SenderPrivateKey       *ecdsa.PrivateKey
	SenderPublicKey        *ecdsa.PublicKey
	SenderWalletAddress    string
	RecipientWalletAddress string
	Value                  float32
	Signature              *Signature
	BroadcastedOn          string
}

type Signature struct {
	R *big.Int
	S *big.Int
}

// This function creates a new transaction object with the given parameters
func (chain *Blockchain) NewTransaction(privateKey string, publicKey string, sender string, recipient string, Value float32) string {
	privateKeyInECDSA, publicKeyInECDSA, err := DecodeStringToECDSA(privateKey, publicKey)
	if err != nil {
		log.Panic(err.Error())
	}

	t := &Transaction{SenderPrivateKey: privateKeyInECDSA, SenderPublicKey: publicKeyInECDSA, SenderWalletAddress: sender, RecipientWalletAddress: recipient, Value: Value}

	t.Signature, err = GenerateSignature(t)
	if err != nil {
		log.Panic(err.Error())
	}
	broadcastedOn := time.Now().Format(time.RFC3339)
	t.BroadcastedOn = broadcastedOn
	t.HashId = generateTransactionHashId(t.Serialize(), broadcastedOn)
	chain.TransactionsQueue = append(chain.TransactionsQueue, t)
	return t.HashId
}

// Generating a signature for the transaction.
func GenerateSignature(transaction *Transaction) (*Signature, error) {
	m, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, transaction.SenderPrivateKey, h[:])
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
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

// Generate Transaction Hash ID
func generateTransactionHashId(bytes []byte, broadcastedOn string) string {
	hash_1_2 := append(bytes[:], broadcastedOn[:]...)
	hash_id := sha256.Sum256(hash_1_2)
	tID := hex.EncodeToString(hash_id[:])
	return tID
}

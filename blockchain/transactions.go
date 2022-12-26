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

	"github.com/rwiteshbera/Hyperion/wallet"
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

const (
	MINING_SNEDER = "THE BLOCKCHAIN"
	MINING_REWARD = 0.05 // MINING_REWARD = 0.05 -> 5%
	MEMPOOL_SIZE  = 2
)

type Miner struct {
	Rewards map[*wallet.Wallet]float32 // Storing the rewards here to send miner once the block is generated
}

var minerRewards = Miner{Rewards: make(map[*wallet.Wallet]float32, 0)}

// This function creates a new transaction object with the given parameters
func (chain *Blockchain) NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, Value float32) *Transaction {
	t := &Transaction{SenderPrivateKey: privateKey, SenderPublicKey: publicKey, SenderWalletAddress: sender, RecipientWalletAddress: recipient, Value: Value}
	t.Signature = GenerateSignature(t)
	chain.TransactionQueue = append(chain.TransactionQueue, t)
	return t
}

// Transfer Balance
func (transaction *Transaction) Transfer(sender *wallet.Wallet, recipient *wallet.Wallet) (bool, float32) {
	if sender.Balance >= transaction.Value {
		// Calculating Gass Fees that will be sent to miner as reward
		gasFees := float32(transaction.Value) * float32(MINING_REWARD)

		// Transfering value to recipent
		recipient.Balance += (transaction.Value - gasFees)

		// Deducting value from sender
		sender.Balance -= transaction.Value

		fmt.Println("Transaction Successful")
		return true, gasFees
	}

	fmt.Printf("Not enough balance : ")
	return false, 0
}

// Generating a signature for the transaction.
func GenerateSignature(transaction *Transaction) *Signature {
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
func (chain *Blockchain) Mining(sender *wallet.Wallet, recipient *wallet.Wallet, miner *wallet.Wallet) {
	// Check if the transactions queue is empty or not?
	if len(chain.TransactionQueue) == 0 {
		return
	}
	// Take first transaction for validation and also remove it from queue
	var tr1 *Transaction
	var err error
	tr1, chain.TransactionQueue, err = dequeue(chain.TransactionQueue)
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

	// If the transaction is valid, add it to mempool
	if valid {
		isTransferred, gasFees := tr1.Transfer(sender, recipient)
		minerRewards.Rewards[miner] += gasFees // Storing the reward value that will be sent once the block is generated

		if isTransferred {
			chain.Mempool = enqueue(chain.Mempool, tr1)

			if len(chain.Mempool) == TransactionsToStoreInBlock {
				// Adding transactions to block
				chain.AddBlock(chain.Mempool, miner.WalletAddress)
				miner.Balance += minerRewards.Rewards[miner] // Transferring rewards to miner that was temporarily stored in minerRewards instance
				minerRewards.Rewards[miner] = 0              // As the balance is transfered, therefore make it 0 (Preventing double spending)

				// Removing the first two transactions from the mempool as it is already added to block
				chain.Mempool = chain.Mempool[TransactionsToStoreInBlock:]
				return
			}
		} else {
			fmt.Println("Transaction Failed")
		}
	} else {
		fmt.Println("Transaction Failed")
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
	err := binary.Write(&buffer, binary.BigEndian, privKeyBytes)
	if err != nil {
		return nil
	}

	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(t.SenderPublicKey)
	err = binary.Write(&buffer, binary.BigEndian, pubKeyBytes)
	if err != nil {
		return nil
	}

	err = binary.Write(&buffer, binary.BigEndian, []byte(t.SenderWalletAddress))
	if err != nil {
		return nil
	}

	err = binary.Write(&buffer, binary.BigEndian, []byte(t.RecipientWalletAddress))
	if err != nil {
		return nil
	}

	err = binary.Write(&buffer, binary.BigEndian, t.Value)
	if err != nil {
		return nil
	}

	return buffer.Bytes()
}

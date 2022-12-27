package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/rwiteshbera/Hyperion/wallet"
)

const (
	MiningRewards              = 3 // It will be given to miner's wallet after successful block creation
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

type Miner struct {
	Wallet  *wallet.Wallet
	Rewards map[*wallet.Wallet]float32 // Storing the rewards here to send miner once the block is generated
}

var minerRewards = Miner{Rewards: make(map[*wallet.Wallet]float32, 0)}
var minerWallet = Miner{Wallet: wallet.CreateWallet(0)}

// This function creates a new transaction object with the given parameters
func (chain *Blockchain) NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, Value float32) *Transaction {
	t := &Transaction{SenderPrivateKey: privateKey, SenderPublicKey: publicKey, SenderWalletAddress: sender, RecipientWalletAddress: recipient, Value: Value}
	t.Signature = GenerateSignature(t)
	chain.TransactionsQueue = append(chain.TransactionsQueue, t)
	return t
}

// Transfer Balance
func (transaction *Transaction) Transfer(sender *wallet.Wallet, recipient *wallet.Wallet) (bool, float32) {
	if sender.Balance >= transaction.Value {
		// Calculating Gas Fees that will be sent to miner as reward
		gasFees := transaction.Value * 0.05

		// Transferring value to recipient
		recipient.Balance += transaction.Value - gasFees

		// Deducting value from sender
		sender.Balance -= transaction.Value
		return true, gasFees
	}

	fmt.Printf("Not enough balance : ")
	return false, 0
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
func (chain *Blockchain) Mining(sender *wallet.Wallet, recipient *wallet.Wallet) {
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

	// If the transaction is valid, add it to mem-pool
	if valid {

		isTransferred, gasFees := tr1.Transfer(sender, recipient)
		minerRewards.Rewards[minerWallet.Wallet] += gasFees // Storing the reward value that will be sent once the block is generated

		if isTransferred {
			chain.Mempool = enqueue(chain.Mempool, tr1)
			fmt.Printf("Transaction Successful : [Balance Transferred : %f] [Gas Fees : %f]\n", tr1.Value, gasFees)
			if len(chain.Mempool) == TransactionsToStoreInBlock {
				// Adding transactions to block
				chain.AddBlock(chain.Mempool, minerWallet.Wallet.GetWalletAddress())
				minerWallet.Wallet.Balance += minerRewards.Rewards[minerWallet.Wallet] + MiningRewards // Transferring rewards to miner that was temporarily stored in minerRewards instance + Extra MiningReward
				minerRewards.Rewards[minerWallet.Wallet] = 0                                           // As the balance is transferred, therefore make it 0 (Preventing double spending)

				// Removing the first two transactions from the mem-pool as it is already added to block
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

// Converting the signature to a string.
func (signature *Signature) Signature() string {
	return fmt.Sprintf("%x%x", signature.R.Bytes(), signature.S.Bytes())
}

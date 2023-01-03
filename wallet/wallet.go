package wallet

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/rwiteshbera/Hyperion/blockchain"
)

/*
	Create New Wallet

Return: type: Wallet 'struct'
*/
func CreateWallet() *Wallet {
	// Generate ECDSA New Private Key (32bytes) - Public-Key (64bytes) Pair
	private, public := generateNewKeyPair()

	// Perform SHA-256 on Public-Key twice
	publicKeyHash := publicKeyHash(public.X.Bytes(), public.Y.Bytes())

	// Add version byte in front of Public-Key-Hash (0x00 for Mainnet)
	versionedHash := append([]byte{version}, publicKeyHash...)

	// Perform SHA-256 twice on Versioned-Hash and take only first 4 bytes
	checkSum := checkSum(versionedHash)

	// Add 4 Checksum bytes after Versioned-Hash
	fullHash := append(versionedHash, checkSum...)

	// Encode byte string to base58
	walletAddress := base58Encode(fullHash)

	w := Wallet{PrivateKey: &private, PublicKey: &public, WalletAddress: string(walletAddress)}

	return &w
}

/*
Get Wallet PrivateKey and PublicKey in string format
*/
func (w *Wallet) GetKeyValuePair() (string, string, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(w.PrivateKey)
	if err != nil {
		return "", "", err
	}
	privateKeyString := base64.StdEncoding.EncodeToString(privateKeyBytes)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(w.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyString := base64.StdEncoding.EncodeToString(publicKeyBytes)

	return privateKeyString, publicKeyString, nil
}

/*
	Get Wallet Address in string format

Return: Wallet-Address (string)
*/
func (w *Wallet) GetWalletAddress() string {
	return w.WalletAddress
}

// A function that returns the balance of the wallet.
func (w *Wallet) GetWalletBalance(chain *blockchain.Blockchain) (float32, error) {
	balances := make(map[string]float32)

	for _, block := range chain.Blocks {
		balances = UpdateWalletBalance(block.TransactionsInBlock, balances)
	}

	balance, ok := balances[w.GetWalletAddress()]
	if !ok {
		return 0, fmt.Errorf("wallet not found")
	}

	return balance, nil
}

// Update Wallet Balance
func UpdateWalletBalance(transactionsInBlock []*blockchain.Transaction, balances map[string]float32) map[string]float32 {
	for _, transaction := range transactionsInBlock {
		balances[transaction.SenderWalletAddress] -= transaction.Value
		balances[transaction.RecipientWalletAddress] += transaction.Value
	}

	return balances
}

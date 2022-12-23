package wallet

import (
	"crypto/ecdsa"
	"fmt"
)

/*
	Create New Wallet

Return: type: Wallet 'struct'
*/
func CreateWallet(b float32) *Wallet {
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
	wallet_addr := base58Encode(fullHash)

	w := Wallet{PrivateKey: &private, PublicKey: &public, WalletAddress: string(wallet_addr), Balance: b}

	return &w
}

/*
	Get ECDSA New Private Key (32bytes) - Public-Key (64bytes) Pair

Return : Private-Key (ecdsa.PrivateKey), Public-Key(ecdsa.PublicKey)
*/
func (w *Wallet) GetKeyValuePair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	return w.PrivateKey, w.PublicKey
}

/*
	Get Private-Key, Public-Key in string format

Return: Private-Key (string), Public Key (string)
*/
func (w *Wallet) GetKeyValuePairInString() (string, string) {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes()), fmt.Sprintf("%x%x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

/*
	Get Wallet Address in string format

Return: Wallet-Address (string)
*/

// A function that returns the balance of the wallet.
func (w *Wallet) GetBalance() float32 {
	return w.Balance
}

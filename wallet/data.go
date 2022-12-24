package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"log"
	"os"
)

const filename = "wallets.bin" // filename + extension

// Add '/' after adding the path
const path = "./wallet/temp/"

var key = []byte("my-encryption-secret-key") // will be added to .env later on with stronger key

type Wallets struct {
	Wallets map[string]*Wallet
}

func InitWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.loadFile()
	return &wallets, err
}

func (wallets *Wallets) AddWallet() string {
	wallet := CreateWallet(0)
	wallets.Wallets[wallet.WalletAddress] = wallet
	return wallet.WalletAddress
}

// LoadWallets() loads wallets file and returns a pointer to a Wallets struct
func LoadWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.loadFile()

	return &wallets, err
}

// Returning a wallet from the wallets map.
func (w *Wallets) GetWallet(address string) Wallet {
	return *w.Wallets[address]
}

// Getting all the addresses from the wallets.
func (w *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range w.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// Loading the wallets from the file.
func (w *Wallets) loadFile() error {
	filepath := path + filename

	contents, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	decrypted, err1 := decryptWalletData(key, contents)
	if err1 != nil {
		log.Panic(err1)
	}

	err = json.Unmarshal(decrypted, &w)
	if err != nil {
		return err
	}

	return nil
}

// Saving the wallets to a file.
func (w *Wallets) SaveFile() {
	filepath := path + filename

	m, err := json.Marshal(w)
	if err != nil {
		log.Panic(err)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Panic(err)
	}

	encrypted, err := encryptWalletData(key, m)

	if err != nil {
		log.Panic("ERROR : ", err.Error())
	}

	err = os.WriteFile(filepath, encrypted, 0644)
	if err != nil {
		log.Panic(err)
	}
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

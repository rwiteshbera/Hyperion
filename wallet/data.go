package wallet

import (
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

func (wallets *Wallets) AddWallet(wallet *Wallet) string {
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
func (wallets *Wallets) GetWallet(address string) Wallet {
	return *wallets.Wallets[address]
}

// Getting all the addresses from  wallets.
func (wallets *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range wallets.Wallets {
		addresses = append(addresses, address)

	}

	return addresses
}

// Loading the wallets from the file.
func (wallets *Wallets) loadFile() error {
	filepath := path + filename

	contents, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	decrypted, err1 := decryptWalletData(key, contents)
	if err1 != nil {
		log.Panic(err1)
	}

	err = json.Unmarshal(decrypted, &wallets)
	if err != nil {
		return err
	}

	return nil
}

// Saving the wallets to a file.
func (wallets *Wallets) SaveFile() {
	filepath := path + filename

	m, err := json.Marshal(wallets)
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

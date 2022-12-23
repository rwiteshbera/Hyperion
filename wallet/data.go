package wallet

import (
	"encoding/json"
	"log"
	"os"
)

const filename = "wallets.json" // filename + extension

// Add '/' after adding the path
const path = "./wallet/temp/"

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

	err = json.Unmarshal(contents, &w)
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

	err = os.WriteFile(filepath, m, 0644)
	if err != nil {
		log.Panic(err)
	}
}

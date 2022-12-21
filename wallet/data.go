package wallet

const path = "./wallets.data"

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
	wallet := CreateWallet()
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
	return nil
}

// Saving the wallets to a file.
func (w *Wallets) SaveFile() {

}

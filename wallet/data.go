package wallet

import (
	"git.mills.io/prologic/bitcask"
)

const path = "./wallet/temp/"
const filename = "wallets.bin"

// Connect with bitcask to store wallets
func ConnectBitcask() *bitcask.Bitcask {
	database, _ := bitcask.Open(path + filename)
	return database
}

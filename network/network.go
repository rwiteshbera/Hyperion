package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/Hyperion/blockchain"
	"github.com/rwiteshbera/Hyperion/wallet"
	"log"
	"net/http"
)

func Server(port *string) {
	fmt.Println("Initializing Blockchain Server...")
	blockchain.BlockchainInstance = blockchain.InitBlockchain()
	blockchain.BlockchainInstance.StartMining()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Visualize the blockchain
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Blockchain": blockchain.BlockchainInstance.GetBlocks()})
	})

	// Check transaction status
	router.GET("/explore/:transaction", func(c *gin.Context) {
		transactionHash := c.Param("transaction")
		all_transactions := blockchain.BlockchainInstance.GetTransactions()

		for _, transaction := range all_transactions {
			if transactionHash == transaction.HashId {
				c.JSON(http.StatusOK, gin.H{"Status": "Confirmed", "sender": transaction.SenderWalletAddress, "recipient": transaction.RecipientWalletAddress, "Amount": transaction.Value, "Broadcasted On": transaction.BroadcastedOn})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"Status": "Pending", "Transaction Id": transactionHash})
	})

	// Check wallet balance
	router.GET("/wallet/:address", func(c *gin.Context) {
		walletAddress := c.Param("address")
		allBlocks := blockchain.BlockchainInstance.GetBlocks()

		balances := make(map[string]float32)
		for _, block := range allBlocks {
			balances = wallet.UpdateWalletBalance(block.TransactionsInBlock, balances)
		}

		balance, ok := balances[walletAddress]
		if !ok {
			c.JSON(http.StatusOK, gin.H{"Wallet Address": walletAddress, "Balance": 0})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Wallet Address": walletAddress, "Balance": balance})
	})

	// Create new transaction
	router.POST("/new", func(c *gin.Context) {
		var txData struct {
			PrivateKey string  `json:"privateKey"`
			PublicKey  string  `json:"publicKey"`
			Sender     string  `json:"sender"`
			Recipient  string  `json:"recipient"`
			Value      float32 `json:"value"`
		}
		if err := c.ShouldBindJSON(&txData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		transactionId := blockchain.BlockchainInstance.NewTransaction(txData.PrivateKey, txData.PublicKey, txData.Sender, txData.Recipient, txData.Value)
		c.JSON(http.StatusOK, gin.H{"Transaction Id": transactionId})
	})

	err := router.Run(":" + *port)
	if err != nil {
		log.Panic(err.Error())
	}

}

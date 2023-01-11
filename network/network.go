package network

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/Hyperion/blockchain"
	"github.com/rwiteshbera/Hyperion/wallet"
	"log"
	"net"
	"net/http"
)

func Server(port *string) {
	fmt.Println("Initializing Blockchain Server...")
	fmt.Println("Warning: The server will not retain any information once it is turned off. Any data on blocks will be lost once the server is stopped.")
	fmt.Printf("\n\n")
	fmt.Printf("Local: http://localhost:%s\n", *port)

	ip, _ := LocalIP() // retrieve the host machine local IP address
	fmt.Printf("On Your Network: http://%s:%s\n", ip, *port)

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

// LocalIP : the function makes use of the Go standard library package "net" to retrieve the host machine local IP address
func LocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

// isPrivateIP : It is used to check the address if it is a private IP or not.
func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/Hyperion/blockchain"
	"log"
	"net/http"
)

var Bcs *blockchain.Blockchain

func Server(port *string) {
	fmt.Println("Initializing Blockchain Server...")
	Bcs = blockchain.InitBlockchain()
	Bcs.StartMining()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "blockchain server is running!"})
	})

	router.GET("/tx", func(c *gin.Context) {
		transactions := Bcs.GetMempool()
		c.JSON(http.StatusOK, gin.H{"transactions": transactions})
	})

	router.POST("/tx", func(c *gin.Context) {
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
		Bcs.NewTransaction(txData.PrivateKey, txData.PublicKey, txData.Sender, txData.Recipient, txData.Value)
		c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
	})

	err := router.Run(":" + *port)
	if err != nil {
		log.Panic(err.Error())
	}

}

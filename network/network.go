package network

import (
	"github.com/rwiteshbera/Hyperion/blockchain"
	"io"
	"log"
	"net/http"
)

type BlockchainServer struct {
	port string
}

func NewBlockchainServer(port string) *BlockchainServer {
	return &BlockchainServer{port}
}

func (chain *BlockchainServer) PORT() string {
	return chain.port
}

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func (chain *BlockchainServer) Run() {
	bcs := blockchain.InitBlockchain()
	bcs.StartMining()
	http.HandleFunc("/", HelloWorld)

	// Start the server
	log.Println("Starting Blockchain server on 0.0.0.0:" + chain.PORT())
	if err := http.ListenAndServe("0.0.0.0:"+chain.PORT(), nil); err != nil {
		log.Fatal(err)
	}
}

func Server(port *string) {
	p := *port
	chain := NewBlockchainServer(p)
	chain.Run()
}

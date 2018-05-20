package server

import (
	"net/http"
	"log"
	"github.com/ita-sammann/toy-chain/blockchain"
	"encoding/json"
)

var Chain blockchain.Blockchain

func blocksListHandler(w http.ResponseWriter, r *http.Request) {
	blocks, err := json.Marshal(Chain.ListBlocks())
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(blocks)
	if err != nil {
		log.Println(err)
	}
}

func StartServer(addr string) {
	if addr == "" {
		addr = ":1138"
	}

	Chain = blockchain.NewBlockchain()

	http.HandleFunc("/blocks/", blocksListHandler)

	log.Println("Listening on", addr)
	http.ListenAndServe(addr, nil)
}

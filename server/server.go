package server

import (
	"net/http"
	"log"
	"github.com/ita-sammann/toy-chain/blockchain"
	"encoding/json"
)

var Chain blockchain.Blockchain

func blocksListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	blocks, err := json.Marshal(Chain.ListBlocks())
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(blocks)
	if err != nil {
		log.Println(err)
	}
}

func mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reqBody struct {
		Data blockchain.BlockData `json:"data"`
	}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	Chain.AddBlock(reqBody.Data)

	http.Redirect(w, r, "/blocks/", 302)
}

func StartServer(addr string) {
	if addr == "" {
		addr = ":1138"
	}

	Chain = blockchain.NewBlockchain()

	http.HandleFunc("/blocks/", blocksListHandler)
	http.HandleFunc("/mine/", mineBlockHandler)

	log.Println("Listening on", addr)
	http.ListenAndServe(addr, nil)
}

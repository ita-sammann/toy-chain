package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ita-sammann/toy-chain/blockchain"
	"github.com/ita-sammann/toy-chain/p2p"
)

var HTTPServer struct {
	chain *blockchain.Blockchain
}

func blocksListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	blocks, err := json.Marshal(HTTPServer.chain.ListBlocks())
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

	HTTPServer.chain.AddBlock(reqBody.Data)

	p2p.BroadcastChain(HTTPServer.chain)

	http.Redirect(w, r, "/blocks/", 302)
}

// StartHTTPServer starts http server
func StartHTTPServer(chain *blockchain.Blockchain, addr string) {
	if addr == "" {
		addr = ":1138"
	}

	HTTPServer.chain = chain

	http.HandleFunc("/blocks/", blocksListHandler)
	http.HandleFunc("/mine/", mineBlockHandler)

	log.Println("Listening HTTP on", addr)
	http.ListenAndServe(addr, nil)
}

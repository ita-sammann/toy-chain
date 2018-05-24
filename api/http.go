package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ita-sammann/toy-chain/blockchain"
	"github.com/ita-sammann/toy-chain/p2p"
)

func blocksListHandler(w http.ResponseWriter, r *http.Request, chain *blockchain.Blockchain) {
	w.Header().Set("Content-Type", "application/json")
	blocks, err := json.Marshal(chain.ListBlocks())
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(blocks)
	if err != nil {
		log.Println(err)
	}
}

func mineBlockHandler(w http.ResponseWriter, r *http.Request, chain *blockchain.Blockchain) {
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

	chain.AddBlock(reqBody.Data)

	p2p.BroadcastChain(chain)

	http.Redirect(w, r, "/blocks/", 302)
}

// StartHTTPServer starts http server
func StartHTTPServer(chain *blockchain.Blockchain, addr string) {
	if addr == "" {
		addr = "127.0.0.1:1138"
	}

	http.HandleFunc("/blocks/", func(w http.ResponseWriter, r *http.Request) { blocksListHandler(w, r, chain) })
	http.HandleFunc("/mine/", func(w http.ResponseWriter, r *http.Request) { mineBlockHandler(w, r, chain) })

	log.Println("Listening HTTP on", addr)
	http.ListenAndServe(addr, nil)
}

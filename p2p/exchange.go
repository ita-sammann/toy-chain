package p2p

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ita-sammann/toy-chain/blockchain"
)

const (
	MsgTypeBlockchain = "blockchain"
	MsgTypeMessage    = "message"
)

type ExchangePayload struct {
	Type      string             `json:"type"`
	PeerID    int                `json:"peerID"`
	Timestamp time.Time          `json:"timestamp"`
	Blocks    []blockchain.Block `json:"blocks,omitempty"`
	Msg       string             `json:"msg,omitempty"`
}

func NewChainPayload(chain *blockchain.Blockchain) ExchangePayload {
	return ExchangePayload{
		MsgTypeBlockchain,
		127,
		time.Now(),
		chain.ListBlocks(),
		"",
	}
}

func NewMsgPayload(msg string) ExchangePayload {
	return ExchangePayload{
		MsgTypeMessage,
		127,
		time.Now(),
		nil,
		msg,
	}
}

func ReplyMsg(msg string) []byte {
	payload, err := json.Marshal(NewMsgPayload(msg))
	if err != nil {
		log.Println(err)
		return []byte(err.Error())
	}
	return payload
}

func listenConnection(conn *Conn, chain *blockchain.Blockchain) {
	for {
		messageType, p, err := conn.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Printf("connection %d closed", conn.id)
			} else {
				log.Println(err)
			}
			conn.conn.Close()
			delete(connPool, conn.id)
			break
		}

		var payload ExchangePayload
		if err := json.Unmarshal(p, &payload); err != nil {
			log.Println("failed to parse payload", err)
			if err := conn.conn.WriteMessage(messageType, ReplyMsg("rejected: failed to parse payload")); err != nil {
				log.Println(err)
				continue
			}
		}

		if payload.Type == MsgTypeBlockchain {
			if err := chain.ReplaceChain(blockchain.NewBlockchainBlocks(payload.Blocks)); err != nil {
				log.Println("rejected:", err)
				if err := conn.conn.WriteMessage(messageType, ReplyMsg("rejected: "+err.Error())); err != nil {
					log.Println(err)
					continue
				}
			} else {
				log.Println("accepted")
				if err := conn.conn.WriteMessage(messageType, ReplyMsg("accepted")); err != nil {
					log.Println(err)
					continue
				}
			}
		} else if payload.Type == MsgTypeMessage {
			log.Println("Recieved message:", payload.Msg)
		} else {
			log.Println("Bad message type:", payload.Type)
		}

	}
}

func SendChain(conn Conn, chain *blockchain.Blockchain) {
	payload, err := json.Marshal(NewChainPayload(chain))
	if err != nil {
		log.Println(err)
		return
	}

	if err := conn.conn.WriteMessage(1, payload); err != nil {
		log.Println(err)
	}
}

func BroadcastChain(chain *blockchain.Blockchain) {
	payload, err := json.Marshal(NewChainPayload(chain))
	if err != nil {
		log.Println(err)
		return
	}

	for _, conn := range connPool {
		if err := conn.conn.WriteMessage(1, payload); err != nil {
			log.Println(err)
			continue
		}
	}
}

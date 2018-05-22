package p2p

import (
	"log"
	"net/http"

	"net/url"

	"github.com/gorilla/websocket"
	"github.com/ita-sammann/toy-chain/blockchain"
)

var incomingConnPool = make([]*websocket.Conn, 128)
var outgoingConnPool = make([]*websocket.Conn, 128)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("WS server: Accepted connection from", conn.RemoteAddr())
	incomingConnPool = append(incomingConnPool, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("WS server: got message from %s: %s", conn.RemoteAddr(), string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}

	//for {
	//	messageType, r, err := conn.NextReader()
	//	if err != nil {
	//		return
	//	}
	//	log.Printf("WS server: got message from %s: %s", conn.RemoteAddr(), )
	//	w, err := conn.NextWriter(messageType)
	//	if err != nil {
	//		panic(err)
	//	}
	//	if _, err := io.Copy(w, r); err != nil {
	//		panic(err)
	//	}
	//	if err := w.Close(); err != nil {
	//		panic(err)
	//	}
	//}
}

// StartWSServer starts http server
func StartWSServer(chain blockchain.Blockchain, addr string) {
	if addr == "" {
		addr = ":11380"
	}

	http.HandleFunc("/", wsHandler)

	log.Println("Listening WS on", addr)
	http.ListenAndServe(addr, nil)
}

func StartWSClient(chain blockchain.Blockchain, addrs []string) {
	for _, addr := range addrs {
		u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
		log.Printf("WS client: connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}

		outgoingConnPool = append(outgoingConnPool, c)
	}
}

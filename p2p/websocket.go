package p2p

import (
	"log"
	"net/http"

	"net/url"

	"github.com/gorilla/websocket"
	"github.com/ita-sammann/toy-chain/blockchain"
)

type Conn struct {
	conn       *websocket.Conn
	isListened bool
}

var connPool = make([]*Conn, 0, 128)
var ConnChan = make(chan *Conn, 16)

func addConnection(wsConn *websocket.Conn, chain *blockchain.Blockchain) {
	conn := &Conn{wsConn, false}
	connPool = append(connPool, conn)
	ConnChan <- conn
	SendChain(*conn, chain)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request, chain *blockchain.Blockchain) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("WS server: Accepted connection from", conn.RemoteAddr())
	addConnection(conn, chain)

	//for {
	//	messageType, p, err := conn.ReadMessage()
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	log.Printf("WS server: got message from %s, type %d: %s", conn.RemoteAddr(), messageType, string(p))
	//	if err := conn.WriteMessage(messageType, p); err != nil {
	//		log.Println(err)
	//		return
	//	}
	//}
}

// StartWSServer starts http server
func StartWSServer(chain *blockchain.Blockchain, addr string) {
	if addr == "" {
		addr = ":11380"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { wsHandler(w, r, chain) })

	log.Println("Listening WS on", addr)
	http.ListenAndServe(addr, nil)
}

func StartWSClient(chain *blockchain.Blockchain, addrs []string) {
	for _, addr := range addrs {
		if addr == "" {
			continue
		}
		u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
		log.Printf("WS client: connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}

		addConnection(c, chain)
	}
}

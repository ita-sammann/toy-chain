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

func wsHandler(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader, chain *blockchain.Blockchain) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("WS server: Accepted connection from", conn.RemoteAddr())
	addConnection(conn, chain)
}

// StartWSServer starts http server
func StartWSServer(chain *blockchain.Blockchain, addr string) {
	if addr == "" {
		addr = "127.0.0.1:11380"
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { wsHandler(w, r, upgrader, chain) })

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

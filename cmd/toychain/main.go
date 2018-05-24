package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/ita-sammann/toy-chain/api"
	"github.com/ita-sammann/toy-chain/blockchain"
	"github.com/ita-sammann/toy-chain/p2p"
)

func main() {
	var httpAddr, wsAddr, wsPeers string
	flag.StringVar(&httpAddr, "http-addr", "127.0.0.1:1138", "HTTP address to listen")
	flag.StringVar(&wsAddr, "ws-addr", "127.0.0.1:11380", "WebSocket address to listen")
	flag.StringVar(&wsPeers, "ws-peers", "", "WebSocket peers")
	flag.Parse()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	exitChan := make(chan int, 1)

	chain := blockchain.NewBlockchain()

	startNetworking(&chain, httpAddr, wsAddr, wsPeers)

	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGHUP:
				log.Println("Hungup")
				exitChan <- 0

			case syscall.SIGINT:
				log.Println("Interrupted")
				exitChan <- 0

			case syscall.SIGTERM:
				fmt.Println("Force stop")
				exitChan <- 0

			case syscall.SIGQUIT:
				fmt.Println("Stop and core dump")
				exitChan <- 0

			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()
	code := <-exitChan
	os.Exit(code)
}

func startNetworking(chain *blockchain.Blockchain, httpAddr, wsAddr, wsPeers string) {
	go api.StartHTTPServer(chain, httpAddr)
	go p2p.StartWSServer(chain, wsAddr)

	peerAddrs := strings.Split(wsPeers, ",")
	go p2p.StartWSClient(chain, peerAddrs)

	go p2p.StartExchange(chain)
}

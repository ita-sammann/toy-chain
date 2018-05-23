package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/ita-sammann/toy-chain/blockchain"
	"github.com/ita-sammann/toy-chain/p2p"
	"github.com/ita-sammann/toy-chain/server"
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
	exit_chan := make(chan int)

	chain := blockchain.NewBlockchain()

	startNetworking(&chain, httpAddr, wsAddr, wsPeers)

	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGHUP:
				log.Println("Hungup")
				exit_chan <- 0

			case syscall.SIGINT:
				log.Println("Interrupted")
				exit_chan <- 0

			case syscall.SIGTERM:
				fmt.Println("Force stop")
				exit_chan <- 0

			case syscall.SIGQUIT:
				fmt.Println("Stop and core dump")
				exit_chan <- 0

			default:
				fmt.Println("Unknown signal.")
				exit_chan <- 1
			}
		}
	}()
	code := <-exit_chan
	os.Exit(code)
}

func startNetworking(chain *blockchain.Blockchain, httpAddr, wsAddr, wsPeers string) {
	go server.StartHTTPServer(chain, httpAddr)
	go p2p.StartWSServer(chain, wsAddr)

	peerAddrs := strings.Split(wsPeers, ",")
	go p2p.StartWSClient(chain, peerAddrs)

	go p2p.StartExchange(chain)
}

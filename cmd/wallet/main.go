package main

import (
	"flag"
	"github.com/martinsaporiti/blockchain-sample/internal/servers"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP port to listen on")
	gateway := flag.String("gateway", "http://localhost:5000", "Blockchain gateway")
	flag.Parse()

	server := servers.NewWalletServer(uint16(*port), *gateway)
	server.Run()
}

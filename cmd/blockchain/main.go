package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/martinsaporiti/blockchain-sample/internal/config"
	"github.com/martinsaporiti/blockchain-sample/internal/controller"
	"github.com/martinsaporiti/blockchain-sample/internal/servers"
	"github.com/martinsaporiti/blockchain-sample/internal/wallet"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port to listen on")
	flag.Parse()

	miningDifficulty := os.Getenv("MINING_DIFFICULTY")
	if miningDifficulty == "" {
		log.Printf("MINING_DIFFICULTY not set, using default value: %s", "5")
		miningDifficulty = "5"
	}

	md, err := strconv.Atoi(miningDifficulty)
	if err != nil {
		log.Panicf("Invalid mining difficulty: %s", miningDifficulty)
	}

	config := config.Config{
		Port:              uint16(*port),
		BlockchainAddress: wallet.New().BlockchainAddress(),
		MiningDifficulty:  md,
	}

	ctrl := controller.New(config)
	server := servers.NewBlockchainServer(config, ctrl)
	server.Start()
}

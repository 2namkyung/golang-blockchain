package main

import (
	"blockchain/pos"
	"log"
	"net"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// 1. PoW , serve TCP ( Using netcat )
// Using TCP and serve TCP Server , PoS ( netcat )
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// create genesis block
	t := time.Now()
	genesisBlock := pos.Block{}
	genesisBlock = pos.Block{0, 0, "Genesis Block", pos.CalcBlockHashForPoS(genesisBlock), "", "", "Genesis Block", t.String()}
	spew.Dump(genesisBlock)
	pos.Blockchain = append(pos.Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	go func() {
		for candidate := range pos.CandidateBlocks {
			pos.MutexForPoS.Lock()
			pos.TempBlocks = append(pos.TempBlocks, candidate)
			pos.MutexForPoS.Unlock()
		}
	}()

	go func() {
		for {
			pos.PickWinner()
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go pos.Conn(conn)
	}

}

package pos

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// Using TCP and serve TCP Server , PoS ( netcat )
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// create genesis block
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, 0, "Genesis Block", CalcBlockHashForPoS(genesisBlock), "", "", "Genesis Block", t.String()}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	go func() {
		for candidate := range CandidateBlocks {
			MutexForPoS.Lock()
			TempBlocks = append(TempBlocks, candidate)
			MutexForPoS.Unlock()
		}
	}()

	go func() {
		for {
			PickWinner()
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go Conn(conn)
	}

}

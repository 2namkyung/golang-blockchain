package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// bcServer handles incoming concurrent Blocks
var bcServer chan []Block

// HandleConn is for PoW
func HandleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM: ")

	scanner := bufio.NewScanner(conn)

	//take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number : %v", scanner.Text(), err)
				continue
			}

			newBlock := GenerateBlock(Blockchain[len(Blockchain)-1], bpm)

			if IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				ReplaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new BPM: ")

		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for range bcServer {
		spew.Dump(Blockchain)
	}
}

/* Network Main */
// Using TCP and serve TCP Server , PoW ( netcat )
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	// create genesis block
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", "", 0, "0"}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", "localhost:"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	// 매번 요청이 올때마다 새로운 Connection을 생성해야한다
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConn(conn)
	}
}

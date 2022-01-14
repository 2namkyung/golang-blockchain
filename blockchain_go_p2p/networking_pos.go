package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

// HandleConnForPoS is for PoS
func HandleConnForPoS(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	// validator address
	var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	io.WriteString(conn, "Enter token balance : ")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number : %v", scanBalance.Text(), err)
		}

		t := time.Now()
		address = CalculateHashForPoS(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nEnter a new BPM : ")
	scanBPM := bufio.NewScanner(conn)

	go func() {
		for {
			// take in BPM from stdin and add it to blockchain after conducting necessary validation
			for scanBPM.Scan() {
				bpm, err := strconv.Atoi(scanBPM.Text())
				// if malicious party tries to mutate the chain with a bad input, delete them as a validator and they lose their staked tokens
				if err != nil {
					log.Printf("%v not a number %v", scanBPM.Text(), err)
					delete(validators, address)
					conn.Close()
				}

				mutexForPoS.Lock()
				oldLastIndex := BlockchainPoS[len(BlockchainPoS)-1]
				mutexForPoS.Unlock()

				// create newBlock for consideration to be forged
				newBlock, err := GenerateBlockForPoS(oldLastIndex, bpm, address)
				if err != nil {
					log.Println(err)
					continue
				}

				if IsBlockValidForPoS(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new BPM : ")
			}
		}
	}()

	// simulate receiving broadcast
	for {
		time.Sleep(time.Minute)
		mutexForPoS.Lock()
		output, err := json.Marshal(BlockchainPoS)
		mutexForPoS.Unlock()

		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(conn, string(output)+"\n")
	}
}

//Using TCP and serve TCP Server , PoS ( netcat )
// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// create genesis block
// 	t := time.Now()
// 	genesisBlock := BlockForPoS{}
// 	genesisBlock = BlockForPoS{0, t.String(), 0, CalculateBlockHashForPoS(genesisBlock), "", ""}
// 	spew.Dump(genesisBlock)
// 	BlockchainPoS = append(BlockchainPoS, genesisBlock)

// 	// start TCP and serve TCP server
// 	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer server.Close()

// 	go func() {
// 		for candidate := range candidateBlocks {
// 			mutexForPoS.Lock()
// 			tempBlocks = append(tempBlocks, candidate)
// 			mutexForPoS.Unlock()
// 		}
// 	}()

// 	go func() {
// 		for {
// 			PickWinner()
// 		}
// 	}()

// 	for {
// 		conn, err := server.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		go HandleConnForPoS(conn)
// 	}

// }

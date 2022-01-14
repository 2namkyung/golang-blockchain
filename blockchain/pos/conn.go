package pos

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

func Conn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	var address string
	io.WriteString(conn, "Enter Token Balance : ")
	scanBalance := bufio.NewScanner(conn)

	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number : %v", scanBalance.Text(), err)
		}

		t := time.Now()
		address = CalcHashForPoS(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nHow much Staking ? : ")
	scanStaking := bufio.NewScanner(conn)

	go func() {
		for {
			for scanStaking.Scan() {
				staking, err := strconv.Atoi(scanStaking.Text())
				if err != nil {
					log.Printf("%v not a number %v", scanStaking.Text(), err)
					delete(validators, address)
					conn.Close()
				}

				MutexForPoS.Lock()
				lastBlock := Blockchain[len(Blockchain)-1]
				MutexForPoS.Unlock()

				newBlock, err := GenerateBlockForPoS(lastBlock, staking, address)
				if err != nil {
					log.Println(err)
					continue
				}

				if IsBlockValidForPoS(newBlock, lastBlock) {
					CandidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nHow much Staking ? : ")
			}
		}
	}()

	// receives broadcasting
	for {
		time.Sleep(time.Minute)
		MutexForPoS.Lock()
		output, err := json.MarshalIndent(Blockchain, "", "")
		MutexForPoS.Unlock()

		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(conn, string(output)+"\n")
	}
}

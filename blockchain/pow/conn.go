package pow

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var BCserver chan []Block

func Conn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a Name : ")

	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			name := scanner.Text()
			_, exists := UserList[name]
			if !exists {
				io.WriteString(conn, "Wait... Making a Wallet\n")
				wallet := MakingWallet()
				UserWallet[name] = wallet
				UserList[name] = 0
			}

			io.WriteString(conn, "What do you want ?")
			scanner.Scan()
			tx := scanner.Text()

			newBlock := GenerateBlock(Blockchain[len(Blockchain)-1], name, tx)

			if IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				ReplaceChain(newBlockchain)
			}

			BCserver <- Blockchain
			io.WriteString(conn, "Enter a Name : ")

		}
	}()

	// receives broadcasting
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.MarshalIndent(Blockchain, "", "")
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output)+"\n")
		}
	}()

	for range BCserver {
		// pretty printer
		spew.Dump(Blockchain)
	}
}

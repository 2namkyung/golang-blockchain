package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var UserList map[string]int
var UserWallet map[string]*Wallet

type Block struct {
	Index       int
	PrevHash    string
	Hash        string
	Difficulty  int
	Nonce       string
	Transaction string
	Timestamp   string
}

var Blockchain []Block
var mutex = &sync.Mutex{}

// Init is the function that makes UserList(map) and genesisBlock
func Init() {
	UserList = make(map[string]int)
	UserWallet = make(map[string]*Wallet)

	// Make genesisBlock
	t := time.Now()
	genesisBlock := Block{0, "", "", 0, "", "Genesis Block", t.String()}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
}

func CalcHash(block Block) string {
	sum := block.PrevHash + block.Nonce + block.Timestamp + string(block.Index)
	h := sha256.New() // SHA256 해시 인스턴스 생성

	h.Write([]byte(sum))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalcHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

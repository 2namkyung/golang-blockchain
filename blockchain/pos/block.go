package pos

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

type Block struct {
	Index       int
	Staking     int
	Validator   string
	PrevHash    string
	Hash        string
	Nonce       string
	Transaction string
	Timestamp   string
}

var Blockchain []Block
var TempBlocks []Block

// handles incoming blocks for validation
var CandidateBlocks = make(chan Block)

// broadcast winning validator to all nodes
var announcements = make(chan string)

var mutex = &sync.Mutex{}
var MutexForPoS = &sync.Mutex{}

var validators = make(map[string]int)

func CalcHashForPoS(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CalcBlockHashForPoS(block Block) string {
	sum := string(block.Index) + block.Timestamp + block.PrevHash
	return CalcHashForPoS(sum)
}

func GenerateBlockForPoS(oldBlock Block, staking int, address string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Staking = staking
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Validator = address
	newBlock.Hash = CalcBlockHashForPoS(newBlock)

	return newBlock, nil
}

func IsBlockValidForPoS(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalcBlockHashForPoS(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func PickWinner() {
	time.Sleep(30 * time.Second)
	mutex.Lock()
	temp := TempBlocks
	mutex.Unlock()

	pool := []string{}
	if len(temp) > 0 {
	OUTER:
		for _, block := range temp {
			for _, node := range pool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// lock list of validators to prevent data race
			MutexForPoS.Lock()
			setValidators := validators
			MutexForPoS.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					pool = append(pool, block.Validator)
				}
			}
		}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		poolWinner := pool[r.Intn(len(pool))]

		// add block to blockchain and let all the other nodes know
		for _, block := range temp {
			if block.Validator == poolWinner {
				MutexForPoS.Lock()
				Blockchain = append(Blockchain, block)
				MutexForPoS.Unlock()
				for range validators {
					announcements <- "\nwinning validator : " + poolWinner
				}
				break
			}
		}
	}

	MutexForPoS.Lock()
	TempBlocks = []Block{}
	MutexForPoS.Unlock()
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

type BlockForPoS struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	Prevhash  string
	Validator string
}

// Blockchain is a series of validated Blocks
var BlockchainPoS []BlockForPoS
var tempBlocks []BlockForPoS

// candidatesBlocks handles incoming blocks for validation
var candidateBlocks = make(chan BlockForPoS)

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)

var mutexForPoS = &sync.Mutex{}

// validators keeps track of open validators and balances
// 각 노드가 가지고 있는 토큰의 수를 나타낸다
var validators = make(map[string]int)

// CalculateHashForPos is a simple SHA256 hasing function
func CalculateHashForPoS(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// CaculateBlockHashForPoS returns the has of all block information
func CalculateBlockHashForPoS(block BlockForPoS) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.Prevhash
	return CalculateHashForPoS(record)
}

func GenerateBlockForPoS(oldBlock BlockForPoS, BPM int, address string) (BlockForPoS, error) {

	var newBlock BlockForPoS

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.Prevhash = oldBlock.Hash
	newBlock.Hash = CalculateBlockHashForPoS(newBlock)
	newBlock.Validator = address

	return newBlock, nil
}

func IsBlockValidForPoS(newBlock, oldBlock BlockForPoS) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.Prevhash {
		return false
	}

	if CalculateBlockHashForPoS(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func PickWinner() {
	time.Sleep(30 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {
		// slightly modified traditional proof of stake algorithm
		// from all validators who submitted a block, weight them by the number of staked tokens
		// in traditional proof of stake, validators can participate without submitting a block to be forged
	OUTER:
		for _, block := range temp {
			// if already in lottery pool, skip
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// lock list of validators to prevent data race
			mutexForPoS.Lock()
			setValidators := validators
			mutexForPoS.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// randomly pick winner from lottery pool
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinnder := lotteryPool[r.Intn(len(lotteryPool))]

		// add block fo winner to blockchain and let all the other nodes know
		for _, block := range temp {
			if block.Validator == lotteryWinnder {
				mutexForPoS.Lock()
				BlockchainPoS = append(BlockchainPoS, block)
				mutexForPoS.Unlock()
				for range validators {
					announcements <- "\nwinning validator : " + lotteryWinnder
				}
				break
			}
		}
	}

	mutexForPoS.Lock()
	tempBlocks = []BlockForPoS{}
	mutexForPoS.Unlock()
}

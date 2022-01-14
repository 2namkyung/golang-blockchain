package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type Message struct {
	BPM int
}

var Blockchain []Block
var mutex = &sync.Mutex{}

type Block struct {
	Index      int
	Timestamp  string
	BPM        int // beats per minute , pulse rate
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

func CalculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash + string(block.Nonce)
	h := sha256.New() // SHA256 해시 인스턴스 생성

	h.Write([]byte(record)) // 해시 인스턴스에 데이터 추가
	hashed := h.Sum(nil)    // 해시 인스턴스에 저장된 데이터의 SHA 해시 값 추출
	return hex.EncodeToString(hashed)
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// 현재 가진 체인보다 새롭게 생성된 체인쪽의 길이가 더 길다면 replace
func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

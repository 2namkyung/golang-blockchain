package pow

import (
	"fmt"
	"strings"
	"time"
)

// set pow difficulty
const difficulty = 1

func GenerateBlock(oldBlock Block, name, tx string) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Transaction = tx
	newBlock.Difficulty = difficulty
	newBlock.Timestamp = t.String()

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex

		hashed := CalcHash(newBlock)
		if !IsHashValid(hashed, newBlock.Difficulty) {
			fmt.Println(name, " : ", hashed, " need to work more !!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println("-------------------------")
			UserList[name]++
			fmt.Println(hashed, " Good !!!")
			fmt.Println(name, "'s Mining Block Count = ", UserList[name])
			newBlock.Hash = hashed
			break
		}
	}

	return newBlock
}

func IsHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

package main

import (
	"fmt"
	"sync"
)

func routineTest(i, count int) {
	fmt.Println(i, count)
}

func main() {

	var wait sync.WaitGroup
	wait.Add(10)

	for i := 0; i < 10; i++ {
		go routineTest(i, 5)
	}
	wait.Wait()
}

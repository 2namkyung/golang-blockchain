package main

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func main() {

	err := qrcode.WriteFile("https://www.naver.com", qrcode.Medium, 256, "test.png")
	if err != nil {
		fmt.Println(err)
	}

}

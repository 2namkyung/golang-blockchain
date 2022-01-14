package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"learngo/keyGen/rsaCrypto"
	"learngo/keyGen/rsaKeyGen"
	"os"
)

func main() {
	reader := rand.Reader
	bitSize := 512

	key, err := rsa.GenerateKey(reader, bitSize)
	rsaKeyGen.CheckError(err)

	publicKey := key.PublicKey

	if _, err := os.Stat("./rsaCrypto/privateKey.pem"); os.IsNotExist(err) {
		fmt.Println("Creating Key Pair . . . ")
		rsaKeyGen.SavePemKey("./rsaCrypto/privateKey.pem", key)
		rsaKeyGen.SavePublicPEMKey("./rsaCrypto/publicKey.pem", publicKey)
	}

	ciphertext, err := rsaCrypto.RSAEncrypt([]byte("blockchain"))
	Check(err)
	fmt.Println("CipherText : ", ciphertext)

	plaintext, err := rsaCrypto.RSADecrypt(ciphertext)
	Check(err)
	fmt.Println("PlainText : ", string(plaintext))
}

func Check(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}

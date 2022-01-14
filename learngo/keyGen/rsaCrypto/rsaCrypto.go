package rsaCrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		return
	}
}

func RSAEncrypt(data []byte) ([]byte, error) {
	publicKey, err := ioutil.ReadFile("./rsaCrypto/publicKey.pem")
	CheckError(err)

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public Key Error")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func RSADecrypt(ciphertext []byte) ([]byte, error) {
	privateKey, err := ioutil.ReadFile("./rsaCrypto/privateKey.pem")
	CheckError(err)

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private Key Error")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

package rsaKeyGen

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}

func SaveGobKey(filename string, key interface{}) {
	file, err := os.Create(filename)
	CheckError(err)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(key)
	CheckError(err)
}

func SavePemKey(filename string, key *rsa.PrivateKey) {
	file, err := os.Create(filename)
	CheckError(err)
	defer file.Close()

	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(file, privateKey)
	CheckError(err)
}

func SavePublicPEMKey(filename string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	CheckError(err)

	pemkey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(filename)
	CheckError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	CheckError(err)
}

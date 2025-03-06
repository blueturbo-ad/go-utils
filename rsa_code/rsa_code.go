package rascode

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func RsaEncrypt(pubKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, fmt.Errorf("publickey error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ParsePKIXPublicKey error %v", err)
	}
	keySize, srcSize := pubInterface.(*rsa.PublicKey).Size(), len(origData)
	pub := pubInterface.(*rsa.PublicKey)
	offSet, onces := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + onces
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData[offSet:endIndex])
		if err != nil {
			return nil, fmt.Errorf("encry once error %v", err)
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}
func RsaDecrypt(pkey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(pkey)
	if block == nil {
		return nil, fmt.Errorf("private key error")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ParsePKCS1PrivateKey error %v", err)
	}
	private := priv.(*rsa.PrivateKey)
	keySize, srcSize := private.Size(), len(ciphertext)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, private, ciphertext[offSet:endIndex])
		if err != nil {
			return nil, fmt.Errorf("decode rsa once error %v", err)
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}

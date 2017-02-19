package main

import (
	"crypto/aes"
	"crypto/cipher"
)

var secret []byte

var globalCipher cipher.Block

func init() {
	globalCipher, _ = aes.NewCipher(secret)
}

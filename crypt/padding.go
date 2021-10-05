package crypt

import (
	"bytes"
	"crypto/aes"
)

func PKCS7Padding(b []byte) []byte {
	padSize := aes.BlockSize - (len(b) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(b, pad...)
}

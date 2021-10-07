package crypt

import "crypto/rand"

func Xor(buf []byte) []byte {
	length := len(buf) / 2
	newBuf := make([]byte, length)
	for i := 0; i < length; i++ {
		newBuf[i] = buf[i] ^ buf[length+i]
	}
	return newBuf
}
func genRandomBytes(length int) []byte {
	buf := make([]byte, length)
	rand.Read(buf)
	return buf
}

func genRandom32Bytes() *[32]byte {
	b := new([32]byte)
	_, _ = rand.Read(b[:])
	return b
}

func genRandom16Bytes() *[16]byte {
	b := new([16]byte)
	_, _ = rand.Read(b[:])
	return b
}

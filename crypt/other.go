package crypt

func Xor(buf []byte) []byte {
	length := len(buf) / 2
	newBuf := make([]byte, length)
	for i := 0; i < length; i++ {
		newBuf[i] = buf[i] ^ buf[length+i]
	}
	return newBuf
}

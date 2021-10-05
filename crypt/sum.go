package crypt

import "crypto/sha256"

func Sha256Sum(attrs ...[]byte) []byte {
	sha := sha256.New()
	for _, attr := range attrs {
		_, err := sha.Write(attr)
		if err != nil {
			return nil
		}
	}
	return sha.Sum(nil)
}

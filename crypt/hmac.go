package crypt

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
)

func SignHmacSha256(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

func SignHmacSha1(key, msg []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

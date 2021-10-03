package line

import (
	"crypto/hmac"
	"crypto/sha1"
)

func SignHmacSha1(key, msg []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

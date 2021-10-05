package crypt

import (
	"crypto/aes"
	"crypto/cipher"
)

func EncryptAesEcb(key, src []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	text := make([]byte, len(src))
	ci.Encrypt(text, src)
	return text, err
}

func DecryptAesEcb(key, src []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	text := make([]byte, len(src))
	ci.Decrypt(text, src)
	return text, nil
}

func EncryptAesGcm(key, iv, add, msg []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(ci)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nil, iv, msg, add), nil
}

func DecryptAesGcm(key, iv, add, msg []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(ci)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, iv, msg, add)
}

func EncryptAesCbc(key, iv, src []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := PKCS7Padding(src)
	text := make([]byte, len(padded))

	cbc := cipher.NewCBCEncrypter(ci, iv)
	cbc.CryptBlocks(text, padded)
	return text, err
}

func DecryptAesCbc(key, iv, src []byte) ([]byte, error) {
	ci, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(ci, iv)
	text := make([]byte, len(src))
	cbc.CryptBlocks(text, src)
	return text, nil
}

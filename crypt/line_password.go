package crypt

import (
	"encoding/base64"
	"github.com/line-api/model/go/model"
)

func EncryptPassword(password string, eekr *model.ExchangeEncryptionKeyResponse, curve25519 *KeyPairForCurve25519) (string, error) {
	eekrPubKey, err := base64.StdEncoding.DecodeString(eekr.PublicKey)
	if err != nil {
		return "", err
	}
	eekrNonce, err := base64.StdEncoding.DecodeString(eekr.Nonce)
	if err != nil {
		return "", err
	}
	secret, err := Curve25519GenSharedSecret(curve25519.PrivateKey[:], eekrPubKey)
	if err != nil {
		return "", err
	}
	masterKey := Sha256Sum([]byte("master_key"), secret, curve25519.Nonce[:], eekrNonce)
	aesKey := Sha256Sum([]byte("aes_key"), masterKey)
	hmacKey := Sha256Sum([]byte("hmac_key"), masterKey)
	enc, err := EncryptAesCbc(aesKey[0:16], aesKey[16:32], []byte(password))
	if err != nil {
		return "", err
	}
	enc = append(enc, SignHmacSha256(hmacKey, enc)...)
	return base64.StdEncoding.EncodeToString(enc), nil
}

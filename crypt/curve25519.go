package crypt

import (
	"encoding/base64"
	"golang.org/x/crypto/curve25519"
)

type KeyPairForCurve25519 struct {
	PrivateKey []byte `json:"private_key"`
	PublicKey  []byte `json:"public_key"`
	Nonce      []byte `json:"nonce"`
}

func generateCurve25519KeyPair() ([]byte, []byte) {
	pub, pri := new([32]byte), genRandom32Bytes()
	curve25519.ScalarBaseMult(pub, pri)
	return pub[:], pri[:]
}

func NewKeyPairForCurve25519() *KeyPairForCurve25519 {
	keyPair := &KeyPairForCurve25519{}
	keyPair.PublicKey, keyPair.PrivateKey = generateCurve25519KeyPair()
	keyPair.Nonce = genRandomBytes(16)
	return keyPair
}

func Curve25519GenSharedSecret(privKey, pubKey []byte) ([]byte, error) {
	return curve25519.X25519(privKey, pubKey)
}

func (k *KeyPairForCurve25519) StringPubKey() string {
	return string(k.PublicKey)
}
func (k *KeyPairForCurve25519) StringPrivKey() string {
	return string(k.PrivateKey)
}
func (k *KeyPairForCurve25519) StringNonce() string {
	return string(k.Nonce)
}

func (k *KeyPairForCurve25519) B64EncodedStringPubKey() string {
	return base64.StdEncoding.EncodeToString(k.PublicKey)
}
func (k *KeyPairForCurve25519) B64EncodedStringPrivKey() string {
	return base64.StdEncoding.EncodeToString(k.PrivateKey)
}
func (k *KeyPairForCurve25519) B64EncodedStringNonce() string {
	return base64.StdEncoding.EncodeToString(k.Nonce)
}

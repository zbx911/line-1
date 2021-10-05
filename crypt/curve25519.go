package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/curve25519"
)

type KeyPairForCurve25519 struct {
	PrivateKey *[32]byte `json:"private_key"`
	PublicKey  *[32]byte `json:"public_key"`
	Nonce      *[16]byte `json:"nonce"`
}

func generateCurve25519KeyPair() (*[32]byte, *[32]byte) {
	pub, pri := new([32]byte), genRandom32Bytes()
	curve25519.ScalarBaseMult(pub, pri)
	return pub, pri
}

func NewKeyPairForCurve25519() *KeyPairForCurve25519 {
	keyPair := &KeyPairForCurve25519{}
	keyPair.PublicKey, keyPair.PrivateKey = generateCurve25519KeyPair()
	keyPair.Nonce = genRandom16Bytes()
	return keyPair
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

func Curve25519GenSharedSecret(privKey, pubKey []byte) ([]byte, error) {
	return curve25519.X25519(privKey, pubKey)
}

func (k *KeyPairForCurve25519) StringPubKey() string {
	return string(k.PublicKey[:])
}
func (k *KeyPairForCurve25519) StringPrivKey() string {
	return string(k.PrivateKey[:])
}
func (k *KeyPairForCurve25519) StringNonce() string {
	return string(k.Nonce[:])
}

func (k *KeyPairForCurve25519) B64EncodedStringPubKey() string {
	return base64.StdEncoding.EncodeToString(k.PublicKey[:])
}
func (k *KeyPairForCurve25519) B64EncodedStringPrivKey() string {
	return base64.StdEncoding.EncodeToString(k.PrivateKey[:])
}
func (k *KeyPairForCurve25519) B64EncodedStringNonce() string {
	return base64.StdEncoding.EncodeToString(k.Nonce[:])
}

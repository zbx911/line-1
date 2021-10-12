package crypt

type E2EESpecVersion int

const (
	E2EESpecVersionUNKNOWN E2EESpecVersion = -1
	E2EESpecVersionV1      E2EESpecVersion = 1
	E2EESpecVersionV2      E2EESpecVersion = 2
)

type KeyStore interface {
	Get(string, int32) (*E2EEKeyPair, bool)
	GetByMid(string) (*E2EEKeyPair, bool)
	Set(string, int32, *E2EEKeyPair)
}

type E2EEKeyPair struct {
	KeyId int32

	PrivateKeyId int32
	PrivateKey   []byte
	PublicKeyId  int32
	PublicKey    []byte
	Version      E2EESpecVersion
}

type E2EEStatus struct {
	SpecVersion E2EESpecVersion
	KeyId       int32
}

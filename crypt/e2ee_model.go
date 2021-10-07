package crypt

import "fmt"

type E2EESpecVersion int

const (
	E2EESpecVersionUNKNOWN E2EESpecVersion = -1
	E2EESpecVersionV1      E2EESpecVersion = 1
	E2EESpecVersionV2      E2EESpecVersion = 2
)

type E2EEKeyPair struct {
	Owner string

	KeyId int32

	PrivateKeyId int32
	PrivateKey   []byte `json:"private_key"`
	PublicKeyId  int32
	PublicKey    []byte `json:"public_key"`
}

type E2EEKeyManagerIF interface {
	Get(mid string, keyId int32) (*E2EEKeyPair, bool)
	Set(mid string, keyId int32, keyPair *E2EEKeyPair)
}

type E2EEKeyStore struct {
	Data map[string]*E2EEKeyPair
}

func (s *E2EEKeyStore) formatKey(keyId int32, mid string) string {
	return fmt.Sprintf("%s_%d", mid, keyId)
}

func (s *E2EEKeyStore) Get(mid string, keyId int32) (*E2EEKeyPair, bool) {
	key, ok := s.Data[s.formatKey(keyId, mid)]
	return key, ok
}
func (s *E2EEKeyStore) Set(mid string, keyId int32, key *E2EEKeyPair) {
	s.Data[s.formatKey(keyId, mid)] = key
}

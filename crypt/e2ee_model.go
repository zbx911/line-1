package crypt

type E2EESpecVersion int

const (
	E2EESpecVersionUNKNOWN E2EESpecVersion = -1
	E2EESpecVersionV1      E2EESpecVersion = 1
	E2EESpecVersionV2      E2EESpecVersion = 2
)

type E2EEKeyPair struct {
	Owner string

	PrivateKeyId int32
	PrivateKey   []byte `json:"private_key"`
	PublicKeyId  int32
	PublicKey    []byte `json:"public_key"`
}

type E2EEKeyManagerIF interface {
	Get(keyId int32) (*E2EEKeyPair, bool)
	Set(keyId int32, keyPair *E2EEKeyPair)
}

type E2EEKeyStore struct {
	Data map[int32]*E2EEKeyPair
}

func (s *E2EEKeyStore) Get(keyId int32) (*E2EEKeyPair, bool) {
	key, ok := s.Data[keyId]
	return key, ok
}
func (s *E2EEKeyStore) Set(keyId int32, key *E2EEKeyPair) {
	s.Data[keyId] = key
}

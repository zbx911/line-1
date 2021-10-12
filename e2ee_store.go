package line

import (
	"fmt"
	"github.com/line-api/line/crypt"
)

type E2EEKeyStore struct {
	Data map[string]*crypt.E2EEKeyPair

	Status map[string]*crypt.E2EEStatus
}

func NewE2EEKeyStore() *E2EEKeyStore {
	return &E2EEKeyStore{
		Data:   map[string]*crypt.E2EEKeyPair{},
		Status: map[string]*crypt.E2EEStatus{},
	}
}

func (s *E2EEKeyStore) formatKey(keyId int32, mid string) string {
	return fmt.Sprintf("%s_%d", mid, keyId)
}

func (s *E2EEKeyStore) Get(mid string, keyId int32) (*crypt.E2EEKeyPair, bool) {
	key, ok := s.Data[s.formatKey(keyId, mid)]
	if ok {
		s.Status[mid] = &crypt.E2EEStatus{
			SpecVersion: key.Version,
			KeyId:       keyId,
		}
	}
	return key, ok
}

func (s *E2EEKeyStore) GetByMid(mid string) (*crypt.E2EEKeyPair, bool) {
	status, ok := s.Status[mid]
	if ok {
		return s.Get(mid, status.KeyId)
	}
	return nil, false
}

func (s *E2EEKeyStore) Set(mid string, keyId int32, key *crypt.E2EEKeyPair) {
	s.Data[s.formatKey(keyId, mid)] = key
	s.Status[mid] = &crypt.E2EEStatus{
		SpecVersion: key.Version,
		KeyId:       keyId,
	}
}

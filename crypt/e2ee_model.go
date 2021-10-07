package crypt

import (
	"fmt"
)

type E2EESpecVersion int

const (
	E2EESpecVersionUNKNOWN E2EESpecVersion = -1
	E2EESpecVersionV1      E2EESpecVersion = 1
	E2EESpecVersionV2      E2EESpecVersion = 2
)

type E2EEKeyPair struct {
	KeyId int32

	PrivateKeyId int32
	PrivateKey   []byte `json:"private_key"`
	PublicKeyId  int32
	PublicKey    []byte `json:"public_key"`
	Version      E2EESpecVersion
}

type E2EEStatus struct {
	SpecVersion E2EESpecVersion
	KeyId       int32
}

type E2EEKeyStore struct {
	Data map[string]*E2EEKeyPair

	Status map[string]*E2EEStatus
}

func NewE2EEKeyStore() *E2EEKeyStore {
	return &E2EEKeyStore{
		Data:   map[string]*E2EEKeyPair{},
		Status: map[string]*E2EEStatus{},
	}
}

func (s *E2EEKeyStore) formatKey(keyId int32, mid string) string {
	return fmt.Sprintf("%s_%d", mid, keyId)
}

func (s *E2EEKeyStore) Get(mid string, keyId int32) (*E2EEKeyPair, bool) {
	key, ok := s.Data[s.formatKey(keyId, mid)]
	if ok {
		s.Status[mid] = &E2EEStatus{
			SpecVersion: key.Version,
			KeyId:       keyId,
		}
	}
	return key, ok
}

func (s *E2EEKeyStore) GetByMid(mid string) (*E2EEKeyPair, bool) {
	status, ok := s.Status[mid]
	if ok {
		return s.Get(mid, status.KeyId)
	}
	return nil, false
}

func (s *E2EEKeyStore) Set(mid string, keyId int32, key *E2EEKeyPair) {
	s.Data[s.formatKey(keyId, mid)] = key
	s.Status[mid] = &E2EEStatus{
		SpecVersion: key.Version,
		KeyId:       keyId,
	}
}

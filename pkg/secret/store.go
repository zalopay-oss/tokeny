package secret

import (
	"github.com/pkg/errors"
	"github.com/zalopay-oss/tokeny/pkg/keyvalue"
)

const secretKeyPrefix = "secret:"

type fallbackStore struct {
	kvStore keyvalue.Store
}

func NewDefaultStore(kvStore keyvalue.Store) Store {
	result, err := newKeyringStore()
	if err == nil {
		return result
	}
	return &fallbackStore{kvStore: kvStore}
}

func (s *fallbackStore) Set(key string, value string) error {
	return s.kvStore.Set(s.composeKey(key), value)
}

func (s *fallbackStore) Get(key string) (string, error) {
	result, err := s.kvStore.Get(s.composeKey(key))
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return "", ErrNoSecret
		}
		return "", err
	}
	return result, nil
}

func (s *fallbackStore) Delete(key string) error {
	err := s.kvStore.Delete(s.composeKey(key))
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return ErrNoSecret
		}
		return err
	}
	return nil
}

func (s *fallbackStore) composeKey(key string) string {
	return secretKeyPrefix + key
}

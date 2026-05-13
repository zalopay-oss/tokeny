package secret

import (
	"github.com/99designs/keyring"
)

const serviceName = "tokeny"

type keyringStore struct {
	keyring keyring.Keyring
}

func newKeyringStore() (*keyringStore, error) {
	backends := secureBackends()
	if len(backends) == 0 {
		return nil, keyring.ErrNoAvailImpl
	}

	backend, err := keyring.Open(keyring.Config{
		ServiceName:     serviceName,
		AllowedBackends: backends,
	})
	if err != nil {
		return nil, err
	}

	return &keyringStore{keyring: backend}, nil
}

func secureBackends() []keyring.BackendType {
	availableBackends := keyring.AvailableBackends()
	result := make([]keyring.BackendType, 0, len(availableBackends))
	for _, backend := range availableBackends {
		if backend == keyring.FileBackend {
			continue
		}
		result = append(result, backend)
	}
	return result
}

func (s *keyringStore) Set(key string, value string) error {
	return s.keyring.Set(keyring.Item{
		Key:  key,
		Data: []byte(value),
	})
}

func (s *keyringStore) Get(key string) (string, error) {
	item, err := s.keyring.Get(key)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", ErrNoSecret
		}
		return "", err
	}
	return string(item.Data), nil
}

func (s *keyringStore) Delete(key string) error {
	err := s.keyring.Remove(key)
	if err != nil && err == keyring.ErrKeyNotFound {
		return ErrNoSecret
	}
	return err
}

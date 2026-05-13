package tokeny

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalopay-oss/tokeny/pkg/keyvalue"
	"github.com/zalopay-oss/tokeny/pkg/secret"
)

type inMemoryKVStore struct {
	values map[string]string
}

func newInMemoryKVStore() *inMemoryKVStore {
	return &inMemoryKVStore{values: map[string]string{}}
}

func (s *inMemoryKVStore) Set(key string, value string) error {
	s.values[key] = value
	return nil
}

func (s *inMemoryKVStore) Get(key string) (string, error) {
	value, ok := s.values[key]
	if !ok {
		return "", keyvalue.ErrNoRecord
	}
	return value, nil
}

func (s *inMemoryKVStore) Delete(key string) error {
	delete(s.values, key)
	return nil
}

func (s *inMemoryKVStore) GetAllWithPrefixed(keyPrefix string) ([]keyvalue.KeyValue, error) {
	result := make([]keyvalue.KeyValue, 0)
	for key, value := range s.values {
		if !strings.HasPrefix(key, keyPrefix) {
			continue
		}
		result = append(result, keyvalue.KeyValue{
			Key:   key,
			Value: value,
		})
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result, nil
}

type inMemorySecretStore struct {
	values map[string]string
}

func newInMemorySecretStore() *inMemorySecretStore {
	return &inMemorySecretStore{values: map[string]string{}}
}

func (s *inMemorySecretStore) Set(key string, value string) error {
	s.values[key] = value
	return nil
}

func (s *inMemorySecretStore) Get(key string) (string, error) {
	value, ok := s.values[key]
	if !ok {
		return "", secret.ErrNoSecret
	}
	return value, nil
}

func (s *inMemorySecretStore) Delete(key string) error {
	delete(s.values, key)
	return nil
}

func TestRepositoryAddStoresSecretOutsideEntryMetadata(t *testing.T) {
	kvStore := newInMemoryKVStore()
	secretStore := newInMemorySecretStore()

	err := NewRepository(kvStore, secretStore).Add("github", " JBSW Y3DP EHPK3PXP ")

	assert.NoError(t, err)
	assert.Equal(t, entryMetadataValue, kvStore.values["entry:github"])
	assert.Equal(t, "JBSWY3DPEHPK3PXP", secretStore.values["entry:github"])
}

func TestRepositoryGenerateMigratesLegacyEntry(t *testing.T) {
	kvStore := newInMemoryKVStore()
	secretStore := newInMemorySecretStore()
	kvStore.values["entry:github"] = "JBSWY3DPEHPK3PXP"

	_, err := NewRepository(kvStore, secretStore).Generate("github")

	assert.NoError(t, err)
	assert.Equal(t, entryMetadataValue, kvStore.values["entry:github"])
	assert.Equal(t, "JBSWY3DPEHPK3PXP", secretStore.values["entry:github"])
	assert.Equal(t, "github", kvStore.values[lastValidKey])
}

func TestRepositoryDeleteRemovesStoredSecret(t *testing.T) {
	kvStore := newInMemoryKVStore()
	secretStore := newInMemorySecretStore()
	kvStore.values["entry:github"] = entryMetadataValue
	kvStore.values[lastValidKey] = "github"
	secretStore.values["entry:github"] = "JBSWY3DPEHPK3PXP"

	err := NewRepository(kvStore, secretStore).Delete("github")

	assert.NoError(t, err)
	_, hasEntry := kvStore.values["entry:github"]
	assert.False(t, hasEntry)
	_, hasSecret := secretStore.values["entry:github"]
	assert.False(t, hasSecret)
	_, hasLastValid := kvStore.values[lastValidKey]
	assert.False(t, hasLastValid)
}

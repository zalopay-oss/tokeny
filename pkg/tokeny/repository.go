package tokeny

import (
	"errors"
	"strings"
	"unicode"

	"github.com/zalopay-oss/tokeny/pkg/keyvalue"
	secretstore "github.com/zalopay-oss/tokeny/pkg/secret"
	"github.com/zalopay-oss/tokeny/pkg/totp"
)

const (
	entryKeyPrefix     = "entry:"
	entryMetadataValue = "__tokeny_secret_ref__"
	lastValidKey       = "last_valid"
)

type repository struct {
	kvStore     keyvalue.Store
	secretStore secretstore.Store
}

func NewRepository(kvStore keyvalue.Store, secretStore secretstore.Store) *repository {
	return &repository{kvStore: kvStore, secretStore: secretStore}
}

func (r *repository) Add(alias string, secret string) error {
	key := r.composeEntryKey(alias)
	_, err := r.kvStore.Get(key)
	if err == nil {
		return ErrEntryExistedBefore
	}
	if !errors.Is(err, keyvalue.ErrNoRecord) {
		return err
	}
	secret = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, secret)

	err = r.secretStore.Set(key, secret)
	if err != nil {
		return err
	}

	err = r.kvStore.Set(key, entryMetadataValue)
	if err != nil {
		rollbackErr := r.secretStore.Delete(key)
		if rollbackErr != nil && !errors.Is(rollbackErr, secretstore.ErrNoSecret) {
			return rollbackErr
		}
		return err
	}

	return nil
}

func (r *repository) Generate(alias string) (totp.Token, error) {
	key := r.composeEntryKey(alias)
	entryValue, err := r.kvStore.Get(key)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return totp.Token{}, ErrNoEntryFound
		}
		return totp.Token{}, err
	}

	secretValue, err := r.getSecret(key, entryValue)
	if err != nil {
		if errors.Is(err, secretstore.ErrNoSecret) {
			return totp.Token{}, ErrNoEntryFound
		}
		return totp.Token{}, err
	}

	g, err := totp.NewGenerator(secretValue)
	if err != nil {
		return totp.Token{}, err
	}
	result := g.Generate()
	err = r.rememberLastValidEntry(alias)
	if err != nil {
		return totp.Token{}, err
	}
	return result, nil
}

func (r *repository) Delete(alias string) error {
	key := r.composeEntryKey(alias)
	entryValue, err := r.kvStore.Get(key)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return ErrNoEntryFound
		}
		return err
	}

	if entryValue == entryMetadataValue {
		err = r.secretStore.Delete(key)
		if err != nil && !errors.Is(err, secretstore.ErrNoSecret) {
			return err
		}
	}

	err = r.kvStore.Delete(key)
	if err != nil {
		return err
	}
	err = r.removeLastValidIfEqual(alias)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) List() ([]string, error) {
	kvs, err := r.kvStore.GetAllWithPrefixed(entryKeyPrefix)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(kvs))
	for i, kv := range kvs {
		result[i] = strings.TrimPrefix(kv.Key, entryKeyPrefix)
	}
	return result, nil
}

func (r *repository) removeLastValidIfEqual(alias string) error {
	lastValid, err := r.kvStore.Get(lastValidKey)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return nil
		}
		return err
	}
	if alias != lastValid {
		return nil
	}
	return r.kvStore.Delete(lastValidKey)
}

func (r *repository) rememberLastValidEntry(alias string) error {
	return r.kvStore.Set(lastValidKey, alias)
}

func (r *repository) getSecret(key string, entryValue string) (string, error) {
	if entryValue != entryMetadataValue {
		err := r.secretStore.Set(key, entryValue)
		if err != nil {
			return "", err
		}

		err = r.kvStore.Set(key, entryMetadataValue)
		if err != nil {
			return "", err
		}

		return entryValue, nil
	}

	return r.secretStore.Get(key)
}

func (r *repository) LastValidEntry() (string, error) {
	result, err := r.kvStore.Get(lastValidKey)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return "", ErrNoEntryFound
		}
		return "", err
	}
	return result, nil
}

func (r *repository) composeEntryKey(alias string) string {
	return entryKeyPrefix + alias
}

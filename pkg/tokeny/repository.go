package tokeny

import (
	"errors"
	"github.com/ltpquang/tokeny/pkg/keyvalue"
	"github.com/ltpquang/tokeny/pkg/totp"
	"strings"
	"unicode"
)

const (
	entryKeyPrefix = "entry:"
	lastValidKey   = "last_valid"
)

type repository struct {
	kvStore keyvalue.Store
}

func NewRepository(kvStore keyvalue.Store) *repository {
	return &repository{kvStore: kvStore}
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
	return r.kvStore.Set(key, secret)
}

func (r *repository) Generate(alias string) (totp.Token, error) {
	key := r.composeEntryKey(alias)
	secret, err := r.kvStore.Get(key)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return totp.Token{}, ErrNoEntryFound
		}
		return totp.Token{}, err
	}
	g, err := totp.NewGenerator(secret)
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
	_, err := r.kvStore.Get(key)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return ErrNoEntryFound
		}
		return err
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
	result := make([]string, len(kvs), len(kvs))
	for i, kv := range kvs {
		result[i] = strings.TrimPrefix(kv.Key, entryKeyPrefix)
	}
	return result, nil
}

func (r *repository) removeLastValidIfEqual(alias string) error {
	lastValid, err := r.LastValidEntry()
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return nil
		}
		return err
	}
	if !(alias != lastValid) {
		return nil
	}
	return r.kvStore.Delete(lastValid)
}

func (r *repository) rememberLastValidEntry(alias string) error {
	return r.kvStore.Set(lastValidKey, alias)
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

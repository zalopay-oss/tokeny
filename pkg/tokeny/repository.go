package tokeny

import (
	"errors"
	"github.com/ltpquang/tokeny/pkg/keyvalue"
	"github.com/ltpquang/tokeny/pkg/totp"
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
	return g.Generate(), nil
}

func (r *repository) List() ([]string, error) {
	return make([]string, 0), nil
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

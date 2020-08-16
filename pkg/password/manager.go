package password

import (
	"encoding/hex"
	"github.com/ltpquang/tokeny/pkg/keyvalue"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	keyPassword = "password"
)

type manager struct {
	kvStore keyvalue.Store
}

func NewManager(kvStore keyvalue.Store) *manager {
	return &manager{kvStore: kvStore}
}

func (m *manager) IsRegistered() (bool, error) {
	_, err := m.kvStore.Get(keyPassword)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, keyvalue.ErrNoRecord) {
		return false, nil
	}
	return false, err
}

func (m *manager) Register(pwd string, rePwd string) error {
	saltedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MaxCost)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(saltedPwd, []byte(rePwd))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordsMismatch
		}
		return err
	}

	err = m.kvStore.Set(keyPassword, hex.EncodeToString(saltedPwd))
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) Login(pwd string) error {
	savedPwdHexStr, err := m.kvStore.Get(keyPassword)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return ErrNotRegistered
		}
		return err
	}

	savedPwd, err := hex.DecodeString(savedPwdHexStr)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(savedPwd, []byte(pwd))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrWrongPassword
		}
		return err
	}

	return nil
}


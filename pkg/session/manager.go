package session

import (
	"errors"
	"fmt"
	"github.com/ltpquang/tokeny/pkg/keyvalue"
	"strconv"
	"time"
)

const (
	keySessionPrefix   = "session:"
	sessionDurationSec = 300
)

type manager struct {
	kvStore keyvalue.Store
}

func NewManager(kvStore keyvalue.Store) *manager {
	result := &manager{kvStore: kvStore}
	err := result.cleanUp()
	if err != nil {
		println(err.Error())
	}
	return result
}

func (m *manager) IsSessionValid(sessionKey string) (bool, error) {
	key := m.composeSessionKey(sessionKey)
	expiredTSStr, err := m.kvStore.Get(key)
	if err != nil {
		if errors.Is(err, keyvalue.ErrNoRecord) {
			return false, nil
		}
		return false, err
	}
	expired, err := m.isTimeStringExpired(expiredTSStr)
	if err != nil {
		return false, err
	}
	return !expired, nil
}

func (m *manager) NewSession(sessionKey string) error {
	key := m.composeSessionKey(sessionKey)
	expiredTS := time.Now().Unix() + sessionDurationSec
	return m.kvStore.Set(key, fmt.Sprintf("%d", expiredTS))
}

func (m *manager) composeSessionKey(sessionKey string) string {
	return keySessionPrefix + sessionKey
}

func (m *manager) cleanUp() error {
	kvs, err := m.kvStore.GetAllWithPrefixed(keySessionPrefix)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		expired, err := m.isTimeStringExpired(kv.Value)
		if err != nil {
			return err
		}
		if !expired {
			continue
		}
		println("delete", kv.Key)
		err = m.kvStore.Delete(kv.Key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) isTimeStringExpired(tsStr string) (bool, error) {
	expiredTS, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return false, err
	}
	return expiredTS < time.Now().Unix(), nil
}

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
	//TODO clean old session data
	return &manager{kvStore: kvStore}
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
	expiredTS, err := strconv.ParseInt(expiredTSStr, 10, 64)
	if err != nil {
		return false, err
	}
	return time.Now().Unix() < expiredTS, nil
}

func (m *manager) NewSession(sessionKey string) error {
	key := m.composeSessionKey(sessionKey)
	expiredTS := time.Now().Unix() + sessionDurationSec
	return m.kvStore.Set(key, fmt.Sprintf("%d", expiredTS))
}

func (m *manager) composeSessionKey(sessionKey string) string {
	return keySessionPrefix + sessionKey
}

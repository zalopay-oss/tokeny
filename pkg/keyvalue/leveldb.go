package keyvalue

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type levelDB struct {
	db *leveldb.DB
}

func NewLevelDBStore(db *leveldb.DB) *levelDB {
	return &levelDB{db: db}
}

func (l *levelDB) Set(key string, value string) error {
	return l.db.Put([]byte(key), []byte(value), nil)
}

func (l *levelDB) Get(key string) (string, error) {
	r, err := l.db.Get([]byte(key), nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return "", ErrNoRecord
	}
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (l *levelDB) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}

func (l *levelDB) GetAllWithPrefixed(keyPrefix string) ([]KeyValue, error) {
	result := make([]KeyValue, 0)
	iter := l.db.NewIterator(util.BytesPrefix([]byte(keyPrefix)), nil)
	for iter.Next() {
		result = append(result, KeyValue{
			Key:   string(iter.Key()),
			Value: string(iter.Value()),
		})
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return result, nil
}

package keyvalue

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	_createSQL      = "CREATE TABLE IF NOT EXISTS kvs (k TEXT PRIMARY KEY, v TEXT)"
	_selectSQL      = "SELECT v FROM kvs WHERE k = ?"
	_deleteSQL      = "DELETE FROM kvs WHERE k = ?"
	_insertSQL      = "INSERT INTO kvs (k, v) VALUES (?, ?)"
	_replaceSQL     = "REPLACE INTO kvs (k, v) VALUES (?, ?)"
	_allPrefixedSQL = "SELECT k, v FROM kvs WHERE k LIKE ?"
)

type sqlStore struct {
	db *sql.DB
}

func NewSQLStore(dbPath string) (*sqlStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	result := &sqlStore{
		db: db,
	}
	return result, result.ensureTable()
}

func (s *sqlStore) ensureTable() error {
	_, err := s.db.Exec(_createSQL)
	return err
}

func (s *sqlStore) Set(key string, value string) error {
	exist, err := s.exist(key)
	if err != nil {
		return err
	}

	q := _insertSQL
	if exist {
		q = _replaceSQL
	}

	_, err = s.db.Exec(q, key, value)

	return err
}

func (s *sqlStore) exist(key string) (bool, error) {
	_, err := s.Get(key)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, ErrNoRecord) {
		return false, nil
	}
	return false, err
}

func (s *sqlStore) Get(key string) (string, error) {
	stmt, err := s.db.Prepare(_selectSQL)

	if err != nil {
		return "", err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	var value string
	err = stmt.QueryRow(key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		}
		return "", err
	}

	return value, nil
}

func (s *sqlStore) Delete(key string) error {
	stmt, err := s.db.Prepare(_selectSQL)
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	if err != nil {
		return err
	}

	var value string
	err = stmt.QueryRow(key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}

	_, err = s.db.Exec(_deleteSQL, key)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) GetAllWithPrefixed(keyPrefix string) ([]KeyValue, error) {
	stmt, err := s.db.Prepare(_allPrefixedSQL)
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	result := make([]KeyValue, 0)

	rows, err := stmt.Query(keyPrefix + "%")
	if err != nil {
		return nil, err
	}

	defer func() {
		if rows.Err() != nil {
			log.Fatalln(rows.Err())
		}
		err := rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for rows.Next() {
		var k, v string
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, err
		}
		result = append(result, KeyValue{Key: k, Value: v})
	}
	return result, nil
}

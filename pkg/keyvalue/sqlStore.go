package keyvalue

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"log"
)

const (
	_createSQL  = "CREATE TABLE IF NOT EXISTS kvs (key TEXT PRIMARY KEY,value TEXT)"
	_selectSQL  = "SELECT value FROM kvs WHERE key = ?"
	_insertSQL  = "INSERT INTO kvs (key, value) VALUES (?, ?)"
	_replaceSQL = "REPLACE INTO kvs (key, value) VALUES (?, ?)"
)

type sqlStore struct {
	dbPath string
}

func NewSQLStore(dbPath string) (*sqlStore, error) {
	result := &sqlStore{dbPath: dbPath}
	return result, result.ensureTable()
}

func (s *sqlStore) openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", s.dbPath)
}

func (s *sqlStore) ensureTable() error {
	db, err := s.openDB()
	if err != nil {
		return err
	}
	defer s.logErr(db.Close())
	_, err = db.Exec(_createSQL)
	return err
}

func (s *sqlStore) logErr(err error) {
	log.Printf("%q", err)
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

	db, err := s.openDB()
	if err != nil {
		return err
	}
	defer s.logErr(db.Close())

	_, err = db.Exec(q, key, value)

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
	db, err := s.openDB()
	if err != nil {
		return "", err
	}
	defer s.logErr(db.Close())

	stmt, err := db.Prepare(_selectSQL)
	if err != nil {
		return "", err
	}
	defer s.logErr(stmt.Close())

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

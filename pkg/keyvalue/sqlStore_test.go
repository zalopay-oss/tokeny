package keyvalue

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type sqlStoreTestSuite struct {
	suite.Suite
	sqlDB *sql.DB
	mock  sqlmock.Sqlmock
}

func (s *sqlStoreTestSuite) SetupTest() {
	conn, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		s.Fail("error creating db connection", err)
	}
	s.sqlDB = conn
	s.mock = mock
}

func (s *sqlStoreTestSuite) TearDownTest() {
	s.mock.ExpectClose()
	err := s.sqlDB.Close()
	if err != nil {
		s.Fail("error closing db connection", err)
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Fail("all expectation were not met", err)
	}
}

func (s *sqlStoreTestSuite) TestGetSuccessfully() {
	key := "foo"
	expected := "bar"
	s.mock.ExpectPrepare(_selectSQL).
		ExpectQuery().
		WithArgs(key).
		WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(expected))

	// Do test
	store := &sqlStore{s.sqlDB}
	result, err := store.Get(key)
	s.NoError(err)
	s.Equal(expected, result, "expect result = %s", expected)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Fail("all expectation were not met", err)
	}

}
func (s *sqlStoreTestSuite) TestGetEmptyKey() {
	key := "foo"
	s.mock.ExpectPrepare(_selectSQL).
		ExpectQuery().
		WithArgs(key).
		WillReturnError(sql.ErrNoRows)

	// Do test
	store := &sqlStore{s.sqlDB}
	result, err := store.Get(key)
	s.EqualError(err, ErrNoRecord.Error())
	s.Empty(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Fail("all expectation were not met", err)
	}

}

func (s *sqlStoreTestSuite) TestDeleteSuccessfully() {
	key := "foo"
	s.mock.ExpectPrepare(_selectSQL).
		ExpectQuery().
		WithArgs(key).
		WillReturnRows(
			sqlmock.NewRows([]string{"k", "v"}).
				AddRow(key, "bar"),
		)
	s.mock.ExpectExec(_deleteSQL).
		WithArgs(key).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Do test
	store := &sqlStore{s.sqlDB}
	err := store.Delete(key)
	s.NoError(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Fail("all expectation were not met", err)
	}

}
func (s *sqlStoreTestSuite) TestDeleteEmptyKey() {
	key := "foo"
	s.mock.ExpectPrepare(_selectSQL).
		ExpectQuery().
		WithArgs(key).
		WillReturnError(sql.ErrNoRows)

	// Do test
	store := &sqlStore{s.sqlDB}
	err := store.Delete(key)
	s.EqualError(err, ErrNoRecord.Error())
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Fail("all expectation were not met", err)
	}

}

func TestSqlStore(t *testing.T) {
	suite.Run(t, new(sqlStoreTestSuite))
}

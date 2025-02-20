package mysqldb

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
)

func MockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	return mockDB, mock, nil
}

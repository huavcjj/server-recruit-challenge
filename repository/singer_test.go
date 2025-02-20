package repository_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pulse227/server-recruit-challenge-sample/infra/mysqldb"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/repository"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SingerRepositorySuite struct {
	mysqldb.DBMYSQLSuite
	singerRepository repository.SingerRepository
}

func TestSingerRepositorySuite(t *testing.T) {
	suite.Run(t, new(SingerRepositorySuite))
}

func (suite *SingerRepositorySuite) SetupSuite() {
	suite.DBMYSQLSuite.SetupSuite()
	suite.singerRepository = repository.NewSingerRepository(suite.DB)
}

func (suite *SingerRepositorySuite) MockDB() sqlmock.Sqlmock {
	mockDB, mock, err := mysqldb.MockDB()
	suite.Require().NoError(err)

	suite.singerRepository = repository.NewSingerRepository(mockDB)
	return mock
}

func (suite *SingerRepositorySuite) AfterTest() {
	suite.singerRepository = repository.NewSingerRepository(suite.DB)
}

func (suite *SingerRepositorySuite) TestSingerRepositoryGetAll() {
	ctx := context.Background()

	singers := []*model.Singer{
		{ID: model.SingerID(1), Name: "Test Singer 1"},
		{ID: model.SingerID(2), Name: "Test Singer 2"},
	}

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, singer := range singers {
		rows.AddRow(singer.ID, singer.Name)
	}

	mock := suite.MockDB()
	mock.ExpectQuery("SELECT id, name FROM singers ORDER BY id").
		WillReturnRows(rows)

	result, err := suite.singerRepository.GetAll(ctx)
	suite.NoError(err)

	suite.Len(result, len(singers))
	for i, singer := range singers {
		suite.Equal(singer.ID, result[i].ID)
		suite.Equal(singer.Name, result[i].Name)
	}

	err = mock.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SingerRepositorySuite) TestSingerRepositoryGet() {
	ctx := context.Background()

	singer := &model.Singer{ID: model.SingerID(1), Name: "Test Singer"}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(singer.ID, singer.Name)

	mock := suite.MockDB()

	mock.ExpectQuery("SELECT id, name FROM singers WHERE id = ?").
		WithArgs(singer.ID).
		WillReturnRows(rows)

	result, err := suite.singerRepository.Get(ctx, singer.ID)
	suite.NoError(err)

	suite.NotNil(result)
	suite.Equal(singer.ID, result.ID)
	suite.Equal(singer.Name, result.Name)

	err = mock.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SingerRepositorySuite) TestSingerRepository_Add() {
	ctx := context.Background()

	singer := &model.Singer{ID: model.SingerID(1), Name: "Test Singer"}

	mock := suite.MockDB()
	mock.ExpectExec("INSERT INTO singers (id, name) VALUES (?, ?)").
		WithArgs(singer.ID, singer.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.singerRepository.Add(ctx, singer)
	suite.NoError(err)

	err = mock.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SingerRepositorySuite) TestSingerRepository_Delete() {
	ctx := context.Background()

	singerID := model.SingerID(1)

	mock := suite.MockDB()
	mock.ExpectExec("DELETE FROM singers WHERE id = ?").
		WithArgs(singerID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.singerRepository.Delete(ctx, singerID)
	suite.NoError(err)

	err = mock.ExpectationsWereMet()
	suite.NoError(err)
}

package repository_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pulse227/server-recruit-challenge-sample/infra/mysqldb"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/repository"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AlbumRepositorySuite struct {
	mysqldb.DBMYSQLSuite
	albumRepository repository.AlbumRepository
}

func TestAlbumRepositorySuite(t *testing.T) {
	suite.Run(t, new(AlbumRepositorySuite))
}

func (suite *AlbumRepositorySuite) SetupSuite() {
	suite.DBMYSQLSuite.SetupSuite()
	suite.albumRepository = repository.NewAlbumRepository(suite.DB)
}

func (suite *AlbumRepositorySuite) MockDB() sqlmock.Sqlmock {
	mockDB, mock, err := mysqldb.MockDB()
	suite.Require().NoError(err)

	suite.albumRepository = repository.NewAlbumRepository(mockDB)
	return mock
}

func (suite *AlbumRepositorySuite) AfterTest() {
	suite.albumRepository = repository.NewAlbumRepository(suite.DB)
}

func (suite *AlbumRepositorySuite) TestAlbumRepositoryAdd() {
	ctx := context.Background()

	album := model.Album{
		ID:       model.AlbumID(1),
		Title:    "Test Album",
		SingerID: model.SingerID(1),
		Singer: &model.Singer{
			ID:   model.SingerID(1),
			Name: "Test Singer",
		},
	}
	mock := suite.MockDB()
	mock.ExpectExec("INSERT INTO albums (id, title, singer_id) VALUES (?, ?, ?)").
		WithArgs(album.ID, album.Title, album.SingerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.albumRepository.Add(ctx, &album)
	suite.Require().NoError(err)

	err = mock.ExpectationsWereMet()
	suite.Require().NoError(err)
}

func (suite *AlbumRepositorySuite) TestAlbumRepositoryGetAll() {
	ctx := context.Background()

	singers := []*model.Singer{
		{
			ID:   model.SingerID(1),
			Name: "Test Singer",
		},
		{
			ID:   model.SingerID(2),
			Name: "Another Singer",
		},
	}
	albums := []*model.Album{
		{
			ID:       model.AlbumID(1),
			Title:    "First Album",
			SingerID: model.SingerID(1),
			Singer:   singers[0],
		},
		{
			ID:       model.AlbumID(2),
			Title:    "Second Album",
			SingerID: model.SingerID(1),
			Singer:   singers[0],
		},
		{
			ID:       model.AlbumID(3),
			Title:    "Third Album",
			SingerID: model.SingerID(2),
			Singer:   singers[1],
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "singer_id", "name"})
	for _, album := range albums {
		rows.AddRow(album.ID, album.Title, album.SingerID, album.Singer.Name)
	}
	mock := suite.MockDB()
	mock.ExpectQuery(
		"SELECT a.id, a.title, a.singer_id, s.name FROM albums a JOIN singers s ON a.singer_id = s.id ORDER BY a.id",
	).WillReturnRows(rows)

	result, err := suite.albumRepository.GetAll(ctx)
	suite.Require().NoError(err)
	suite.Require().Len(result, len(albums))

	for i, album := range albums {
		suite.Assert().Equal(album.ID, result[i].ID)
		suite.Assert().Equal(album.Title, result[i].Title)
		suite.Assert().Equal(album.SingerID, result[i].SingerID)
		suite.Assert().Equal(album.Singer.ID, result[i].Singer.ID)
		suite.Assert().Equal(album.Singer.Name, result[i].Singer.Name)
	}

	err = mock.ExpectationsWereMet()
	suite.Require().NoError(err)
}

func (suite *AlbumRepositorySuite) TestAlbumRepositoryGet() {
	ctx := context.Background()

	album := &model.Album{
		ID:       model.AlbumID(1),
		Title:    "First Album",
		SingerID: model.SingerID(1),
		Singer: &model.Singer{
			ID:   model.SingerID(1),
			Name: "Test Singer",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "singer_id", "name"}).
		AddRow(album.ID, album.Title, album.SingerID, album.Singer.Name)

	mock := suite.MockDB()
	mock.ExpectQuery(
		"SELECT a.id, a.title, a.singer_id, s.name FROM albums a JOIN singers s ON a.singer_id = s.id WHERE a.id = ?",
	).WithArgs(album.ID).WillReturnRows(rows)

	result, err := suite.albumRepository.Get(ctx, album.ID)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)

	suite.Assert().Equal(model.AlbumID(1), result.ID)
	suite.Assert().Equal("First Album", result.Title)
	suite.Assert().Equal(model.SingerID(1), result.SingerID)
	suite.Assert().Equal(model.SingerID(1), result.Singer.ID)
	suite.Assert().Equal("Test Singer", result.Singer.Name)

	err = mock.ExpectationsWereMet()
	suite.Require().NoError(err)
}

func (suite *AlbumRepositorySuite) TestAlbumRepositoryDelete() {
	ctx := context.Background()
	albumID := model.AlbumID(1)

	mock := suite.MockDB()

	mock.ExpectExec("DELETE FROM albums WHERE id = ?").
		WithArgs(albumID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.albumRepository.Delete(ctx, albumID)
	suite.Require().NoError(err)

	mock.ExpectQuery(
		"SELECT a.id, a.title, a.singer_id, s.name FROM albums a JOIN singers s ON a.singer_id = s.id WHERE a.id = ?",
	).WithArgs(albumID).
		WillReturnError(sql.ErrNoRows)

	result, err := suite.albumRepository.Get(ctx, albumID)
	suite.Require().Error(err)
	suite.Require().Nil(result)
	suite.Require().Equal(repository.ErrorAlbumNotFound, err)

	err = mock.ExpectationsWereMet()
	suite.Require().NoError(err)

}

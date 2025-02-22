package service_test

import (
	"context"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MockAlbumRepository struct {
	mock.Mock
}

func NewMockAlbumRepository() *MockAlbumRepository {
	return &MockAlbumRepository{}
}

func (m *MockAlbumRepository) GetAll(ctx context.Context) ([]*model.Album, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Album), args.Error(1)
}
func (m *MockAlbumRepository) Get(ctx context.Context, id model.AlbumID) (*model.Album, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Album), args.Error(1)
}
func (m *MockAlbumRepository) Add(ctx context.Context, album *model.Album) error {
	args := m.Called(ctx, album)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}
func (m *MockAlbumRepository) Delete(ctx context.Context, id model.AlbumID) error {
	args := m.Called(ctx, id)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

type AlbumServiceSuite struct {
	suite.Suite
	albumService        service.AlbumService
	mockAlbumRepository *MockAlbumRepository
}

func TestAlbumServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumServiceSuite))
}

func (suite *AlbumServiceSuite) SetupSuite() {
	suite.mockAlbumRepository = NewMockAlbumRepository()
	suite.albumService = service.NewAlbumService(suite.mockAlbumRepository)
}

func (suite *AlbumServiceSuite) TestAlbumServiceGetAlbumListService() {
	ctx := context.Background()

	albums := []*model.Album{
		{
			ID:       model.AlbumID(1),
			Title:    "First Album",
			SingerID: model.SingerID(1),
			Singer:   &model.Singer{ID: model.SingerID(1), Name: "Test Singer"},
		},
		{
			ID:       model.AlbumID(2),
			Title:    "Second Album",
			SingerID: model.SingerID(1),
			Singer:   &model.Singer{ID: model.SingerID(1), Name: "Test Singer"},
		},
	}

	suite.mockAlbumRepository.On("GetAll", ctx).Return(albums, nil)

	result, err := suite.albumService.GetAlbumListService(ctx)

	suite.Require().Nil(err)
	suite.Assert().Equal(albums, result)

	suite.mockAlbumRepository.AssertExpectations(suite.T())
}

func (suite *AlbumServiceSuite) TestAlbumServiceGetAlbumService() {
	ctx := context.Background()

	album := &model.Album{
		ID:       model.AlbumID(1),
		Title:    "First Album",
		SingerID: model.SingerID(1),
		Singer:   &model.Singer{ID: model.SingerID(1), Name: "Test Singer"},
	}

	suite.mockAlbumRepository.On("Get", ctx, album.ID).Return(album, nil)

	result, err := suite.albumService.GetAlbumService(ctx, album.ID)

	suite.Require().Nil(err)
	suite.Assert().Equal(album, result)

	suite.mockAlbumRepository.AssertExpectations(suite.T())
}

func (suite *AlbumServiceSuite) TestAlbumServicePostAlbumService() {
	ctx := context.Background()

	album := &model.Album{
		ID:       model.AlbumID(1),
		Title:    "New Album",
		SingerID: model.SingerID(1),
		Singer:   &model.Singer{ID: model.SingerID(1), Name: "Test Singer"},
	}

	suite.mockAlbumRepository.On("Add", ctx, album).Return(nil)

	err := suite.albumService.PostAlbumService(ctx, album)

	suite.Require().Nil(err)

	suite.mockAlbumRepository.AssertExpectations(suite.T())
}

func (suite *AlbumServiceSuite) TestAlbumServiceDeleteAlbumService() {
	ctx := context.Background()
	id := model.AlbumID(1)

	suite.mockAlbumRepository.On("Delete", ctx, id).Return(nil)

	err := suite.albumService.DeleteAlbumService(ctx, id)

	suite.Require().Nil(err)

	suite.mockAlbumRepository.AssertExpectations(suite.T())
}

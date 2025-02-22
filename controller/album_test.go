package controller_test

import (
	"context"
	"encoding/json"
	"github.com/pulse227/server-recruit-challenge-sample/controller"
	"github.com/pulse227/server-recruit-challenge-sample/dto"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockAlbumService struct {
	mock.Mock
}

func NewMockAlbumService() *MockAlbumService {
	return &MockAlbumService{}
}

func (m *MockAlbumService) GetAlbumListService(ctx context.Context) ([]*model.Album, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Album), args.Error(1)
}

func (m *MockAlbumService) GetAlbumService(ctx context.Context, albumID model.AlbumID) (*model.Album, error) {
	args := m.Called(ctx, albumID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Album), args.Error(1)
}

func (m *MockAlbumService) PostAlbumService(ctx context.Context, album *model.Album) error {
	args := m.Called(ctx, album)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

func (m *MockAlbumService) DeleteAlbumService(ctx context.Context, albumID model.AlbumID) error {
	args := m.Called(ctx, albumID)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

type AlbumControllerSuite struct {
	suite.Suite
	albumController  controller.AlbumController
	mockAlbumService *MockAlbumService
}

func TestAlbumControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumControllerSuite))
}
func (suite *AlbumControllerSuite) SetupTest() {
	suite.mockAlbumService = NewMockAlbumService()
	suite.albumController = controller.NewAlbumController(suite.mockAlbumService)
}

func (suite *AlbumControllerSuite) TestGetAlbumsSuccess() {
	req := httptest.NewRequest(http.MethodGet, "/albums", nil)
	rr := httptest.NewRecorder()

	albums := []*model.Album{
		{
			ID:    model.AlbumID(1),
			Title: "Album 1",
			Singer: &model.Singer{
				ID:   model.SingerID(1),
				Name: "Singer 1",
			},
		},
		{
			ID:    model.AlbumID(2),
			Title: "Album 2",
			Singer: &model.Singer{
				ID:   model.SingerID(1),
				Name: "Singer 1",
			},
		},
	}

	suite.mockAlbumService.On("GetAlbumListService", req.Context()).Return(albums, nil)
	suite.albumController.GetAlbums(rr, req)

	suite.Require().Equal(http.StatusOK, rr.Code)

	var res []*dto.AlbumResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Assert().Len(res, 2)
	suite.Assert().Equal(1, res[0].ID)
	suite.Assert().Equal(2, res[1].ID)
	suite.Assert().Equal("Album 1", res[0].Title)
	suite.Assert().Equal("Album 2", res[1].Title)
	suite.Assert().Equal(1, res[0].Singer.ID)
	suite.Assert().Equal(1, res[1].Singer.ID)
	suite.Assert().Equal("Singer 1", res[0].Singer.Name)
	suite.Assert().Equal("Singer 1", res[1].Singer.Name)

	suite.mockAlbumService.AssertExpectations(suite.T())
}

func (suite *AlbumControllerSuite) TestGetAlbumSuccess() {
	//req := httptest.NewRequest(http.MethodGet, "/albums/1", nil)
	//rr := httptest.NewRecorder()
	//
	//album := &model.Album{
	//	ID:       model.AlbumID(1),
	//	Title:    "Test Album",
	//	SingerID: model.SingerID(2),
	//	Singer: &model.Singer{
	//		ID:   model.SingerID(2),
	//		Name: "Test Singer",
	//	},
	//}
	//
	//request := dto.GetAlbumRequest{ID: 1}
	//requestID := request.ToModel()
	//
	//suite.mockAlbumService.On("GetAlbumService", req.Context(), *requestID).Return(album, nil)
	//suite.albumController.GetAlbum(rr, req)
	//
	//suite.Equal(http.StatusOK, rr.Code)
	//
	//var res dto.AlbumResponse
	//err := json.NewDecoder(rr.Body).Decode(&res)
	//suite.NoError(err)
	//
	//suite.Equal(1, res.ID)
	//suite.Equal("Test Album", res.Title)
	//suite.Equal(2, res.Singer.ID)
	//suite.Equal("Test Singer", res.Singer.Name)
	//
	//suite.mockAlbumService.AssertExpectations(suite.T())
}

func (suite *AlbumControllerSuite) TestCreateAlbumSuccess() {
	body := `{"id":1,"title":"New Album","singer_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/albums", strings.NewReader(body))
	rr := httptest.NewRecorder()

	var album dto.CreateAlbumRequest
	err := json.NewDecoder(strings.NewReader(body)).Decode(&album)
	suite.Require().NoError(err)

	suite.mockAlbumService.On("PostAlbumService", req.Context(), album.ToModel()).Return(nil)
	suite.albumController.CreateAlbum(rr, req)

	suite.Require().Equal(http.StatusCreated, rr.Code)

	var res dto.CreateAlbumResponse
	err = json.NewDecoder(rr.Body).Decode(&res)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Assert().Equal(1, res.ID)
	suite.Assert().Equal("New Album", res.Title)
	suite.Assert().Equal(1, res.SingerID)

	suite.mockAlbumService.AssertExpectations(suite.T())
}

func (suite *AlbumControllerSuite) TestDeleteAlbumSuccess() {
	//req := httptest.NewRequest(http.MethodDelete, "/albums/1", nil)
	//rr := httptest.NewRecorder()
	//
	//request := dto.DeleteAlbumRequest{ID: 1}
	//albumID := request.ToModel()
	//
	//suite.mockAlbumService.On("DeleteAlbumService", req.Context(), &albumID).Return(nil)
	//suite.albumController.DeleteAlbum(rr, req)
	//
	//suite.Equal(http.StatusNoContent, rr.Code)
	//
	//suite.mockAlbumService.AssertExpectations(suite.T())
}

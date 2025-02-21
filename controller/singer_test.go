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

type MockSingerService struct {
	mock.Mock
}

func NewMockSingerService() *MockSingerService {
	return &MockSingerService{}
}

func (m *MockSingerService) GetSingerListService(ctx context.Context) ([]*model.Singer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Singer), args.Error(1)
}

func (m *MockSingerService) GetSingerService(ctx context.Context, singerID model.SingerID) (*model.Singer, error) {
	args := m.Called(ctx, singerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Singer), args.Error(1)
}

func (m *MockSingerService) PostSingerService(ctx context.Context, singer *model.Singer) error {
	args := m.Called(ctx, singer)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

func (m *MockSingerService) DeleteSingerService(ctx context.Context, singerID model.SingerID) error {
	args := m.Called(ctx, singerID)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

type SingerControllerSuite struct {
	suite.Suite
	singerController  controller.SingerController
	mockSingerService *MockSingerService
}

func TestSingerControllerSuite(t *testing.T) {
	suite.Run(t, new(SingerControllerSuite))
}

func (suite *SingerControllerSuite) SetupTest() {
	suite.mockSingerService = NewMockSingerService()
	suite.singerController = controller.NewSingerController(suite.mockSingerService)
}

func (suite *SingerControllerSuite) TestGetSingerListHandler() {
	req := httptest.NewRequest(http.MethodGet, "/singers", nil)
	rr := httptest.NewRecorder()

	singers := []*model.Singer{
		{
			ID:   model.SingerID(1),
			Name: "Singer 1",
		},
		{
			ID:   model.SingerID(2),
			Name: "Singer 2",
		},
	}

	suite.mockSingerService.On("GetSingerListService", req.Context()).Return(singers, nil)
	suite.singerController.GetSingerListHandler(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var res []*dto.SingerResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	suite.NoError(err)

	suite.Len(res, 2)
	suite.Equal(1, res[0].ID)
	suite.Equal("Singer 1", res[0].Name)
	suite.Equal(2, res[1].ID)
	suite.Equal("Singer 2", res[1].Name)

	suite.mockSingerService.AssertExpectations(suite.T())
}

func (suite *SingerControllerSuite) TestGetSingerDetailHandler() {}

func (suite *SingerControllerSuite) TestPostSingerHandler() {
	body := `{"id":1,"name":"Singer 1"}`
	req := httptest.NewRequest(http.MethodPost, "/singers", strings.NewReader(body))
	rr := httptest.NewRecorder()

	var singer dto.CreateSingerRequest
	err := json.NewDecoder(strings.NewReader(body)).Decode(&singer)
	suite.NoError(err)

	suite.mockSingerService.On("PostSingerService", req.Context(), singer.ToModel()).Return(nil)
	suite.singerController.PostSingerHandler(rr, req)

	suite.Equal(http.StatusCreated, rr.Code)

	var res dto.SingerResponse
	err = json.NewDecoder(rr.Body).Decode(&res)
	suite.NoError(err)

	suite.Equal(1, res.ID)
	suite.Equal("Singer 1", res.Name)

	suite.mockSingerService.AssertExpectations(suite.T())
}

func (suite *SingerControllerSuite) TestDeleteSingerHandler() {}

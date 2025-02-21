package service_test

import (
	"context"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MockSingerRepository struct {
	mock.Mock
}

func NewMockSingerRepository() *MockSingerRepository {
	return &MockSingerRepository{}
}

func (m *MockSingerRepository) GetAll(ctx context.Context) ([]*model.Singer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Singer), args.Error(1)
}
func (m *MockSingerRepository) Get(ctx context.Context, id model.SingerID) (*model.Singer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Singer), args.Error(1)
}
func (m *MockSingerRepository) Add(ctx context.Context, singer *model.Singer) error {
	args := m.Called(ctx, singer)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}
func (m *MockSingerRepository) Delete(ctx context.Context, id model.SingerID) error {
	args := m.Called(ctx, id)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

type SingerServiceSuite struct {
	suite.Suite
	singerService        service.SingerService
	mockSingerRepository *MockSingerRepository
}

func TestSingerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SingerServiceSuite))
}

func (suite *SingerServiceSuite) SetupSuite() {
	suite.mockSingerRepository = NewMockSingerRepository()
	suite.singerService = service.NewSingerService(suite.mockSingerRepository)
}

func (suite *SingerServiceSuite) TestSingerServiceGetSingerListService() {
	ctx := context.Background()

	singers := []*model.Singer{
		{ID: model.SingerID(1), Name: "Test Singer 1"},
		{ID: model.SingerID(2), Name: "Test Singer 2"},
	}
	suite.mockSingerRepository.On("GetAll", ctx).Return(singers, nil)

	result, err := suite.singerService.GetSingerListService(ctx)
	suite.Assert().Nil(err)
	suite.Assert().Equal(singers, result)
	suite.Assert().Equal(len(singers), len(result))
	suite.Assert().Equal(singers[0].ID, result[0].ID)
	suite.Assert().Equal(singers[0].Name, result[0].Name)
	suite.Assert().Equal(singers[1].ID, result[1].ID)
	suite.Assert().Equal(singers[1].Name, result[1].Name)
	suite.mockSingerRepository.AssertExpectations(suite.T())
}

func (suite *SingerServiceSuite) TestSingerServiceGetSingerService() {
	ctx := context.Background()

	singer := &model.Singer{ID: model.SingerID(1), Name: "Test Singer"}
	suite.mockSingerRepository.On("Get", ctx, singer.ID).Return(singer, nil)

	result, err := suite.singerService.GetSingerService(ctx, singer.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(singer, result)
	suite.mockSingerRepository.AssertExpectations(suite.T())
}

func (suite *SingerServiceSuite) TestSingerServicePostSingerService() {
	ctx := context.Background()

	singer := &model.Singer{ID: model.SingerID(1), Name: "Test Singer"}
	suite.mockSingerRepository.On("Add", ctx, singer).Return(nil)

	err := suite.singerService.PostSingerService(ctx, singer)
	suite.Assert().Nil(err)
	suite.mockSingerRepository.AssertExpectations(suite.T())
}

func (suite *SingerServiceSuite) TestSingerServiceDeleteSingerService() {
	ctx := context.Background()

	id := model.SingerID(1)
	suite.mockSingerRepository.On("Delete", ctx, id).Return(nil)

	err := suite.singerService.DeleteSingerService(ctx, id)
	suite.Assert().Nil(err)
	suite.mockSingerRepository.AssertExpectations(suite.T())
}

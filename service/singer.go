package service

import (
	"context"

	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/repository"
)

type SingerService interface {
	GetSingerListService(ctx context.Context) ([]*model.Singer, error)
	GetSingerService(ctx context.Context, singerID model.SingerID) (*model.Singer, error)
	PostSingerService(ctx context.Context, singer *model.Singer) error
	DeleteSingerService(ctx context.Context, singerID model.SingerID) error
}

type singerService struct {
	singerRepository repository.SingerRepository
}

var _ SingerService = (*singerService)(nil)

func NewSingerService(singerRepository repository.SingerRepository) SingerService {
	return &singerService{singerRepository: singerRepository}
}

func (s *singerService) GetSingerListService(ctx context.Context) ([]*model.Singer, error) {
	return s.singerRepository.GetAll(ctx)
}

func (s *singerService) GetSingerService(ctx context.Context, singerID model.SingerID) (*model.Singer, error) {
	return s.singerRepository.Get(ctx, singerID)
}

func (s *singerService) PostSingerService(ctx context.Context, singer *model.Singer) error {
	if err := singer.Validate(); err != nil {
		return err
	}

	return s.singerRepository.Add(ctx, singer)
}

func (s *singerService) DeleteSingerService(ctx context.Context, singerID model.SingerID) error {
	return s.singerRepository.Delete(ctx, singerID)
}

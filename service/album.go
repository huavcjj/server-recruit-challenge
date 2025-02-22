package service

import (
	"context"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/pulse227/server-recruit-challenge-sample/repository"
)

type AlbumService interface {
	GetAlbumListService(ctx context.Context) ([]*model.Album, error)
	GetAlbumService(ctx context.Context, albumID model.AlbumID) (*model.Album, error)
	PostAlbumService(ctx context.Context, album *model.Album) error
	DeleteAlbumService(ctx context.Context, albumID model.AlbumID) error
}

type albumService struct {
	albumRepository repository.AlbumRepository
}

var _ AlbumService = (*albumService)(nil)

func NewAlbumService(albumRepository repository.AlbumRepository) AlbumService {
	return &albumService{albumRepository: albumRepository}
}

func (s *albumService) GetAlbumListService(ctx context.Context) ([]*model.Album, error) {
	return s.albumRepository.GetAll(ctx)
}

func (s *albumService) GetAlbumService(ctx context.Context, albumID model.AlbumID) (*model.Album, error) {
	return s.albumRepository.Get(ctx, albumID)
}

func (s *albumService) PostAlbumService(ctx context.Context, album *model.Album) error {
	if err := album.Validate(); err != nil {
		return err
	}
	return s.albumRepository.Add(ctx, album)
}

func (s *albumService) DeleteAlbumService(ctx context.Context, albumID model.AlbumID) error {
	return s.albumRepository.Delete(ctx, albumID)
}

package dto

import "github.com/pulse227/server-recruit-challenge-sample/model"

type CreateAlbumRequest struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	SingerID int    `json:"singer_id"`
}

func (r *CreateAlbumRequest) ToModel() *model.Album {
	return &model.Album{
		ID:       model.AlbumID(r.ID),
		Title:    r.Title,
		SingerID: model.SingerID(r.SingerID),
	}
}

type CreateAlbumResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	SingerID int    `json:"singer_id"`
}

func NewCreateAlbumResponse(album *model.Album) *CreateAlbumResponse {
	return &CreateAlbumResponse{
		ID:       int(album.ID),
		Title:    album.Title,
		SingerID: int(album.SingerID),
	}
}

type GetAlbumRequest struct {
	ID int `json:"id"`
}

func (r *GetAlbumRequest) ToModel() *model.AlbumID {
	id := model.AlbumID(r.ID)
	return &id
}

type DeleteAlbumRequest struct {
	ID int `json:"id"`
}

func (r *DeleteAlbumRequest) ToModel() *model.AlbumID {
	id := model.AlbumID(r.ID)
	return &id
}

type AlbumResponse struct {
	ID     int             `json:"id"`
	Title  string          `json:"title"`
	Singer *SingerResponse `json:"singer"`
}

func NewAlbumResponse(album *model.Album) *AlbumResponse {
	return &AlbumResponse{
		ID:    int(album.ID),
		Title: album.Title,
		Singer: &SingerResponse{
			ID:   int(album.Singer.ID),
			Name: album.Singer.Name,
		},
	}
}

func NewAlbumsResponse(albums []*model.Album) []*AlbumResponse {
	responses := make([]*AlbumResponse, len(albums))
	for i, album := range albums {
		responses[i] = NewAlbumResponse(album)
	}
	return responses
}

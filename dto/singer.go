package dto

import "github.com/pulse227/server-recruit-challenge-sample/model"

type SingerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewSingerResponse(singer *model.Singer) *SingerResponse {
	return &SingerResponse{
		ID:   int(singer.ID),
		Name: singer.Name,
	}
}

func NewSingersResponse(singers []*model.Singer) []*SingerResponse {
	res := make([]*SingerResponse, 0)
	for _, singer := range singers {
		res = append(res, NewSingerResponse(singer))
	}
	return res
}

type GetSingerRequest struct {
	ID int `json:"id"`
}

func (r *GetSingerRequest) ToModel() *model.SingerID {
	id := model.SingerID(r.ID)
	return &id
}

type CreateSingerRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (r *CreateSingerRequest) ToModel() *model.Singer {
	return &model.Singer{
		ID:   model.SingerID(r.ID),
		Name: r.Name,
	}
}

type DeleteSingerRequest struct {
	ID int `json:"id"`
}

func (r *DeleteSingerRequest) ToModel() *model.SingerID {
	id := model.SingerID(r.ID)
	return &id
}

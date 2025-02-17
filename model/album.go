package model

type AlbumID int

type Album struct {
	ID       AlbumID  `json:"id"`
	Title    string   `json:"title"`
	SingerID SingerID `json:"singer_id"`
}

func (a *Album) Validate() error {
	if a.Title == "" {
		return ErrInvalidParam
	}
	if len(a.Title) > 255 {
		return ErrInvalidParam
	}
	return nil
}

package repository

import "errors"

var (
	ErrorSingerNotFound = errors.New("singer not found")
	ErrorAlbumNotFound  = errors.New("album not found")
)

package mysqldb

import "errors"

var (
	ErrorSingerNotFound = errors.New("singer not found")
	ErrorAlbumNotFound  = errors.New("album not found")
)

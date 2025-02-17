package mysqldb

import (
	"database/sql"

	"github.com/pulse227/server-recruit-challenge-sample/repository"
)

type albumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) repository.AlbumRepository {
	return &albumRepository{
		db: db,
	}
}

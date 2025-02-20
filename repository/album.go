package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"log/slog"
)

type AlbumRepository interface {
	GetAll(ctx context.Context) ([]*model.Album, error)
	Get(ctx context.Context, id model.AlbumID) (*model.Album, error)
	Add(ctx context.Context, album *model.Album) error
	Delete(ctx context.Context, id model.AlbumID) error
}
type albumRepository struct {
	db *sql.DB
}

var _ AlbumRepository = (*albumRepository)(nil)

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	return &albumRepository{
		db: db,
	}
}

func (r *albumRepository) GetAll(ctx context.Context) ([]*model.Album, error) {
	query := `
		SELECT a.id, a.title, a.singer_id, s.name
		FROM albums a
		JOIN singers s ON a.singer_id = s.id
		ORDER BY a.id 
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	albums := make([]*model.Album, 0)
	for rows.Next() {
		album := model.Album{}
		singer := model.Singer{}
		if err = rows.Scan(&album.ID, &album.Title, &album.SingerID, &singer.Name); err != nil {
			return nil, err
		}
		singer.ID = album.SingerID
		album.Singer = &singer

		albums = append(albums, &album)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return albums, nil
}

func (r *albumRepository) Get(ctx context.Context, id model.AlbumID) (*model.Album, error) {
	query := `
		SELECT a.id, a.title, a.singer_id, s.name
		FROM albums a
		JOIN singers s ON a.singer_id = s.id
		WHERE a.id = ?
	`

	album := model.Album{}
	singer := model.Singer{}
	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&album.ID, &album.Title, &album.SingerID, &singer.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorAlbumNotFound
		}
		return nil, err
	}

	singer.ID = album.SingerID
	album.Singer = &singer

	return &album, nil
}

func (r *albumRepository) Add(ctx context.Context, album *model.Album) error {
	query := `INSERT INTO albums (id, title, singer_id) VALUES (?, ?, ?)`
	if _, err := r.db.ExecContext(ctx, query, album.ID, album.Title, album.SingerID); err != nil {
		return err
	}
	return nil
}

func (r *albumRepository) Delete(ctx context.Context, id model.AlbumID) error {
	query := `DELETE FROM albums WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrorAlbumNotFound
	}

	return nil
}

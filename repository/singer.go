package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/pulse227/server-recruit-challenge-sample/model"
)

type SingerRepository interface {
	GetAll(ctx context.Context) ([]*model.Singer, error)
	Get(ctx context.Context, id model.SingerID) (*model.Singer, error)
	Add(ctx context.Context, singer *model.Singer) error
	Delete(ctx context.Context, id model.SingerID) error
}
type singerRepository struct {
	db *sql.DB
}

var _ SingerRepository = (*singerRepository)(nil)

func NewSingerRepository(db *sql.DB) SingerRepository {
	return &singerRepository{
		db: db,
	}
}
func (r *singerRepository) GetAll(ctx context.Context) ([]*model.Singer, error) {
	query := `SELECT id, name FROM singers ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	singers := make([]*model.Singer, 0)
	for rows.Next() {
		singer := model.Singer{}
		if err = rows.Scan(&singer.ID, &singer.Name); err != nil {
			return nil, err
		}
		singers = append(singers, &singer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return singers, nil
}

func (r *singerRepository) Get(ctx context.Context, id model.SingerID) (*model.Singer, error) {
	query := `SELECT id, name FROM singers WHERE id = ?`
	singer := model.Singer{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(&singer.ID, &singer.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorSingerNotFound
	} else if err != nil {
		return nil, err
	}

	return &singer, nil
}

func (r *singerRepository) Add(ctx context.Context, singer *model.Singer) error {
	query := `INSERT INTO singers (id, name) VALUES (?, ?)`
	if _, err := r.db.ExecContext(ctx, query, singer.ID, singer.Name); err != nil {
		return err
	}
	return nil
}

func (r *singerRepository) Delete(ctx context.Context, id model.SingerID) error {
	query := `DELETE FROM singers WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrorSingerNotFound
	}

	return nil
}

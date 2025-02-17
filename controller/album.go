package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pulse227/server-recruit-challenge-sample/dto"
	"github.com/pulse227/server-recruit-challenge-sample/service"
	"log/slog"
	"net/http"
	"strconv"
)

type AlbumController interface {
	GetAlbums(w http.ResponseWriter, r *http.Request)
	GetAlbum(w http.ResponseWriter, r *http.Request)
	CreateAlbum(w http.ResponseWriter, r *http.Request)
	DeleteAlbum(w http.ResponseWriter, r *http.Request)
}
type albumController struct {
	service service.AlbumService
}

var _ AlbumController = (*albumController)(nil)

func NewAlbumController(s service.AlbumService) AlbumController {
	return &albumController{service: s}
}

// GetAlbums GET /albums
func (a albumController) GetAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := a.service.GetAlbumListService(r.Context())
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := dto.NewAlbumsResponse(albums)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// GetAlbum GET /albums/{id}
func (a albumController) GetAlbum(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	req := dto.GetAlbumRequest{ID: ID}
	albumID := req.ToModel()
	album, err := a.service.GetAlbumService(r.Context(), *albumID)

	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := dto.NewAlbumResponse(album)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// CreateAlbum POST /albums
func (a albumController) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateAlbumRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	album := req.ToModel()
	if err := a.service.PostAlbumService(r.Context(), album); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err = fmt.Errorf("album ID already exists: %w", err)
			errorHandler(w, r, http.StatusConflict, err.Error())
			return
		}
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := dto.NewCreateAlbumResponse(album)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// DeleteAlbum DELETE /albums/{id}
func (a albumController) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	req := dto.DeleteAlbumRequest{ID: ID}
	albumID := req.ToModel()
	if err = a.service.DeleteAlbumService(r.Context(), *albumID); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

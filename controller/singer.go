package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pulse227/server-recruit-challenge-sample/dto"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/pulse227/server-recruit-challenge-sample/service"
)

type SingerController interface {
	GetSingerListHandler(w http.ResponseWriter, r *http.Request)
	GetSingerDetailHandler(w http.ResponseWriter, r *http.Request)
	PostSingerHandler(w http.ResponseWriter, r *http.Request)
	DeleteSingerHandler(w http.ResponseWriter, r *http.Request)
}

type singerController struct {
	service service.SingerService
}

var _ SingerController = (*singerController)(nil)

func NewSingerController(s service.SingerService) SingerController {
	return &singerController{service: s}
}

// GetSingerListHandler GET /singers
func (c *singerController) GetSingerListHandler(w http.ResponseWriter, r *http.Request) {
	singers, err := c.service.GetSingerListService(r.Context())
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := dto.NewSingersResponse(singers)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// GetSingerDetailHandler GET /singers/{id}
func (c *singerController) GetSingerDetailHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	req := dto.GetSingerRequest{ID: ID}
	albumID := req.ToModel()
	singer, err := c.service.GetSingerService(r.Context(), *albumID)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := dto.NewSingerResponse(singer)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// PostSingerHandler POST /singers
func (c *singerController) PostSingerHandler(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateSingerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("invalid body param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	singer := req.ToModel()
	if err := c.service.PostSingerService(r.Context(), singer); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := dto.NewSingerResponse(singer)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// DeleteSingerHandler DELETE /singers/{id}
func (c *singerController) DeleteSingerHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	req := dto.DeleteSingerRequest{ID: ID}
	singerID := req.ToModel()
	if err = c.service.DeleteSingerService(r.Context(), *singerID); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			err = fmt.Errorf("cannot delete singer: related albums exist: %w", err)
			errorHandler(w, r, http.StatusConflict, err.Error())
			return
		}
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

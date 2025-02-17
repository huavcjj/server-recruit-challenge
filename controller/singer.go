package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/pulse227/server-recruit-challenge-sample/model"
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
	if err = json.NewEncoder(w).Encode(singers); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// GetSingerDetailHandler GET /singers/{id}
func (c *singerController) GetSingerDetailHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	singerID, err := strconv.Atoi(idString)
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	singer, err := c.service.GetSingerService(r.Context(), model.SingerID(singerID))
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(singer); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// PostSingerHandler POST /singers
func (c *singerController) PostSingerHandler(w http.ResponseWriter, r *http.Request) {
	var singer *model.Singer
	if err := json.NewDecoder(r.Body).Decode(&singer); err != nil {
		err = fmt.Errorf("invalid body param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.PostSingerService(r.Context(), singer); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(singer); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return
	}
}

// DeleteSingerHandler DELETE /singers/{id}
func (c *singerController) DeleteSingerHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	singerID, err := strconv.Atoi(idString)
	if err != nil {
		err = fmt.Errorf("invalid path param: %w", err)
		errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err = c.service.DeleteSingerService(r.Context(), model.SingerID(singerID)); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

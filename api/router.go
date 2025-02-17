package api

import (
	"net/http"

	"github.com/pulse227/server-recruit-challenge-sample/api/middleware"
	"github.com/pulse227/server-recruit-challenge-sample/controller"
	"github.com/pulse227/server-recruit-challenge-sample/infra/mysqldb"
	"github.com/pulse227/server-recruit-challenge-sample/service"
)

func NewRouter(
	dbUser, dbPass, dbHost, dbName string,
) (http.Handler, error) {

	dbClient, err := mysqldb.Initialize(dbUser, dbPass, dbHost, dbName)
	if err != nil {
		return nil, err
	}
	if err = dbClient.Ping(); err != nil {
		return nil, err
	}

	singerRepo := mysqldb.NewSingerRepository(dbClient)
	singerService := service.NewSingerService(singerRepo)
	singerController := controller.NewSingerController(singerService)

	albumRepo := mysqldb.NewAlbumRepository(dbClient)
	albumService := service.NewAlbumService(albumRepo)
	albumController := controller.NewAlbumController(albumService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /singers", singerController.GetSingerListHandler)
	mux.HandleFunc("GET /singers/{id}", singerController.GetSingerDetailHandler)
	mux.HandleFunc("POST /singers", singerController.PostSingerHandler)
	mux.HandleFunc("DELETE /singers/{id}", singerController.DeleteSingerHandler)

	mux.HandleFunc("GET /albums", albumController.GetAlbums)
	mux.HandleFunc("GET /albums/{id}", albumController.GetAlbum)
	mux.HandleFunc("POST /albums", albumController.CreateAlbum)
	mux.HandleFunc("DELETE /albums/{id}", albumController.DeleteAlbum)

	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux, nil
}

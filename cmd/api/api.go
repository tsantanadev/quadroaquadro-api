package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tsantanadev/social-api/internal/rest"
	"github.com/tsantanadev/social-api/internal/store"
)

type application struct {
	config     config
	store      store.Storage
	tmdbClient rest.TMDBClient
}

type config struct {
	addr       string
	db         dbConfig
	TMDBConfig TMDBConfig
}

type TMDBConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mountRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 5))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthHandler)

		r.Route("/movies", func(r chi.Router) {
			r.Get("/", app.listMoviesHandler)
			r.Post("/", app.createMovieHandler)
		})
	})

	return r
}

func (app *application) Run(mux http.Handler) error {

	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	log.Printf("Starting server on %s", app.config.addr)
	return srv.ListenAndServe()
}

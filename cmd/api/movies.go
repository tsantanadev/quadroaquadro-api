package main

import (
	"net/http"

	"github.com/tsantanadev/social-api/internal/store"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var movie store.Movie
	if err := readJSON(w, r, &movie); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := app.store.Movies.Create(r.Context(), &movie); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, movie)
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := app.store.Movies.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, movies)
}

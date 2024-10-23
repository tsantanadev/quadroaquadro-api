package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tsantanadev/quadroaquadro/internal/mapper"
	"github.com/tsantanadev/quadroaquadro/internal/vo"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var request vo.MovieRequest
	if err := readJSON(w, r, &request); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	movieExists, err := app.store.Movies.Exists(r.Context(), request.Id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if movieExists {
		writeJSONError(w, http.StatusConflict, "movie already exists")
		return
	}

	tmdbResponse, err := app.tmdbClient.GetMovie(request.Id)
	movie := mapper.MapTMDBMovieToStoreMovie(&tmdbResponse, &request)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := app.store.Movies.Create(r.Context(), movie); err != nil {
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

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	movie, err := app.store.Movies.Get(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, movie)
}

package main

import (
	"fmt"
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

	responses := make([]vo.MovieResponse, 0)
	for _, movie := range movies {
		images := app.store.Images.GetImagesByMovieId(movie.ID)
		imagesResponse := make([]vo.Image, 0)
		for _, image := range images {
			url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", app.config.Bucket, image.ID)
			imageResponse := vo.Image{
				Level: image.Level,
				URL:   url,
			}
			imagesResponse = append(imagesResponse, imageResponse)
		}

		movieResponse := vo.MovieResponse{
			ID:            movie.ID,
			Title:         movie.Title,
			Origin:        movie.Origin,
			Tags:          movie.Tags,
			ReleaseDate:   movie.ReleaseDate,
			OriginalTitle: movie.OriginalTitle,
			PosterPath:    movie.PosterPath,
			Genres:        movie.Genres,
			Images:        imagesResponse,
		}

		responses = append(responses, movieResponse)
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, responses)
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

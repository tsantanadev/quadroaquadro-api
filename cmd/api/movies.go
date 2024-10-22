package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/tsantanadev/social-api/internal/rest"
	"github.com/tsantanadev/social-api/internal/store"
)

type movieRequest struct {
	Id   int
	Tags []string
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var request movieRequest
	if err := readJSON(w, r, &request); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	tmdbResponse, err := app.tmdbClient.GetMovie(request.Id)
	movie := store.Movie{
		ID:            request.Id,
		Title:         tmdbResponse.Title,
		OriginalTitle: tmdbResponse.OriginalTitle,
		Origin:        tmdbResponse.OriginCountry,
		Tags:          captalizeTags(&request),
		PosterPath:    tmdbResponse.PosterPath,
		Genres:        extractGenres(tmdbResponse),
		ReleaseDate:   parseDate(tmdbResponse.ReleaseDate),
	}

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

func captalizeTags(request *movieRequest) []string {
	var capitalized []string
	for _, t := range request.Tags {
		capitalized = append(capitalized, strings.ToUpper(t))
	}
	return capitalized
}

func extractGenres(tmdbResponse rest.TMDBMovie) []string {
	var extracted []string
	for _, g := range tmdbResponse.Genres {
		extracted = append(extracted, g.Name)
	}
	return extracted
}

func parseDate(date string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{} // return zero value of time.Time in case of error
	}
	return parsedDate
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := app.store.Movies.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, movies)
}

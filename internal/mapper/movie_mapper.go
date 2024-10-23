package mapper

import (
	"strings"
	"time"

	"github.com/tsantanadev/quadroaquadro/internal/rest"
	"github.com/tsantanadev/quadroaquadro/internal/store"
	"github.com/tsantanadev/quadroaquadro/internal/vo"
)

func MapTMDBMovieToStoreMovie(tmdbMovie *rest.TMDBMovie, request *vo.MovieRequest) store.Movie {
	return store.Movie{
		ID:            request.Id,
		Title:         tmdbMovie.Title,
		OriginalTitle: tmdbMovie.OriginalTitle,
		Origin:        tmdbMovie.OriginCountry,
		Tags:          captalizeTags(request),
		PosterPath:    tmdbMovie.PosterPath,
		Genres:        extractGenres(tmdbMovie),
		ReleaseDate:   parseDate(tmdbMovie.ReleaseDate),
	}
}

func captalizeTags(request *vo.MovieRequest) []string {
	var capitalized []string
	for _, t := range request.Tags {
		capitalized = append(capitalized, strings.ToUpper(t))
	}
	return capitalized
}

func extractGenres(tmdbResponse *rest.TMDBMovie) []string {
	var extracted []string
	for _, g := range tmdbResponse.Genres {
		extracted = append(extracted, g.Name)
	}
	return extracted
}

func parseDate(date string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}
	}
	return parsedDate
}

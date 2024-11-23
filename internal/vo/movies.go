package vo

import (
	"time"
)

type MovieRequest struct {
	Id   int
	Tags []string
}

type MovieResponse struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Origin        []string  `json:"origin"`
	Tags          []string  `json:"tags"`
	ReleaseDate   time.Time `json:"releaseDate"`
	OriginalTitle string    `json:"originalTitle"`
	PosterPath    string    `json:"posterPath"`
	Genres        []string  `json:"genres"`
	Images        []Image   `json:"images"`
}

type Image struct {
	URL   string `json:"url"`
	Level int    `json:"level"`
}

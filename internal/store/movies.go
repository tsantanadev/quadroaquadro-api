package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Origin        []string  `json:"origin"`
	Tags          []string  `json:"tags"`
	ReleaseDate   time.Time `json:"release_date"`
	OriginalTitle string    `json:"original_title"`
	PosterPath    string    `json:"poster_path"`
	Genres        []string  `json:"genres"`
}

type MovieStore struct {
	db *sql.DB
}

func (s *MovieStore) Create(ctx context.Context, movie Movie) error {
	query := `
		INSERT INTO movies (id, title, origin, release_date, tags, original_title, poster_path, genres)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		movie.ID,
		movie.Title,
		pq.Array(movie.Origin),
		movie.ReleaseDate,
		pq.Array(movie.Tags),
		movie.OriginalTitle,
		movie.PosterPath,
		pq.Array(movie.Genres),
	)
	return err
}

func (s *MovieStore) List(ctx context.Context) ([]Movie, error) {
	query := `
		SELECT id, title, origin, tags, release_date, original_title, poster_path, genres 
		FROM movies
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Origin,
			pq.Array(&movie.Tags),
			&movie.ReleaseDate,
			&movie.OriginalTitle,
			&movie.PosterPath,
			pq.Array(&movie.Genres),
		); err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

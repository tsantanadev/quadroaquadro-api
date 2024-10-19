package store

import (
	"context"
	"database/sql"
	"time"
)

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Origin      string    `json:"origin"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	ReleaseDate time.Time `json:"release_date"`
	Status      string    `json:"status"`
}

type MovieStore struct {
	db *sql.DB
}

func (s *MovieStore) Create(ctx context.Context, movie *Movie) error {
	query := `
		INSERT INTO movies (title, origin, category, release_date, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		movie.Title,
		movie.Origin,
		movie.Category,
		movie.ReleaseDate,
		movie.Status,
	).Scan(&movie.ID)
}

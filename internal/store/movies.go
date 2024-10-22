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
	Origin        string    `json:"origin"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	ReleaseDate   time.Time `json:"release_date"`
	Status        string    `json:"status"`
	OriginalTitle string    `json:"original_title"`
}

type MovieStore struct {
	db *sql.DB
}

func (s *MovieStore) Create(ctx context.Context, movie *Movie) error {
	query := `
		INSERT INTO movies (title, origin, category, release_date, status, tags, original_title)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		movie.Title,
		movie.Origin,
		movie.Category,
		movie.ReleaseDate,
		movie.Status,
		pq.Array(movie.Tags),
		movie.OriginalTitle,
	).Scan(&movie.ID)
}

func (s *MovieStore) List(ctx context.Context) ([]Movie, error) {
	query := `
		SELECT id, title, origin, category, release_date, status, tags, original_title
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
			&movie.Category,
			&movie.ReleaseDate,
			&movie.Status,
			pq.Array(&movie.Tags),
			&movie.OriginalTitle,
		); err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

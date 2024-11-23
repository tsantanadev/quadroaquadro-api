package store

import "database/sql"

type Image struct {
	ID      string `json:"id"`
	MovieId int    `json:"movie_id"`
	Level   int    `json:"level"`
}

type ImageStore struct {
	db *sql.DB
}

func (s *ImageStore) Create(image *Image) error {
	query := `
		INSERT INTO images (id, movie_id, level)
		VALUES ($1, $2, $3) RETURNING id
	`

	return s.db.QueryRow(
		query,
		image.ID,
		image.MovieId,
		image.Level,
	).Scan(&image.ID)
}

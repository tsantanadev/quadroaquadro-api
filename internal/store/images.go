package store

import "database/sql"

type Image struct {
	ID       int    `json:"id"`
	Movie_id int    `json:"movie_id"`
	Level    string `json:"level"`
	Url      string `json:"url"`
}

type ImageStore struct {
	db *sql.DB
}

func (s *ImageStore) Create(image *Image) error {
	query := `
		INSERT INTO images (movie_id, level, url)
		VALUES ($1, $2, $3) RETURNING id
	`

	return s.db.QueryRow(
		query,
		image.Movie_id,
		image.Level,
		image.Url,
	).Scan(&image.ID)
}

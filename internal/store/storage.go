package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Movies interface {
		Create(context.Context, Movie) error
		List(context.Context) ([]Movie, error)
		Exists(context.Context, int) (bool, error)
		Get(context.Context, int) (*Movie, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Movies: &MovieStore{db},
		Users:  &UserStore{db},
	}
}

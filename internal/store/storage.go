package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Movies interface {
		Create(context.Context, *Movie) error
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

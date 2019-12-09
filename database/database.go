package database

import (
	users "github.com/cagodoy/tenpo-users-api"
	"github.com/cagodoy/tenpo-users-api/database/postgres"
	"github.com/jmoiron/sqlx"
)

// Store ...
type Store interface {
	UserGet(*users.Query) (*users.User, error)
	UserCreate(*users.User) error

	// TODO(ca): below methods are not implemented
	// List() ([]*users.User, error)
	// Update(*users.User) error
	// Delete(*users.User) error
}

// NewPostgres ...
func NewPostgres(dsn string) (Store, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &postgres.UserStore{
		Store: db,
	}, nil
}

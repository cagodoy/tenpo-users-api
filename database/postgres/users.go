package postgres

import (
	"errors"

	"github.com/Masterminds/squirrel"
	users "github.com/cagodoy/tenpo-users-api"
	"github.com/jmoiron/sqlx"
)

// UserStore ...
type UserStore struct {
	Store *sqlx.DB
}

// UserGet ...
func (us *UserStore) UserGet(q *users.Query) (*users.User, error) {
	query := squirrel.Select("*").From("users").Where("deleted_at is null")

	if q.ID == "" && q.Email == "" {
		return nil, errors.New("must proovide a query")
	}

	if q.ID != "" {
		query = query.Where("id = ?", q.ID)
	}

	if q.Email != "" {
		query = query.Where("email = ?", q.Email)
	}

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := us.Store.QueryRowx(sql, args...)

	c := &users.User{}
	if err := row.StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// UserCreate ...
func (us *UserStore) UserCreate(u *users.User) error {

	sql, args, err := squirrel.
		Insert("users").
		Columns("email", "name", "password").
		Values(u.Email, u.Name, u.Password).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sql, args...)
	if err := row.StructScan(u); err != nil {
		return err
	}

	return nil
}

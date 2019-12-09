package service

import (
	users "github.com/cagodoy/tenpo-users-api"
	"github.com/cagodoy/tenpo-users-api/database"
)

// NewUsers ...
func NewUsers(store database.Store) *Users {
	return &Users{
		Store: store,
	}
}

// Users ...
type Users struct {
	Store database.Store
}

// GetByID ...
func (us *Users) GetByID(id string) (*users.User, error) {
	q := &users.Query{
		ID: id,
	}

	return us.Store.UserGet(q)
}

// GetByEmail ...
func (us *Users) GetByEmail(email string) (*users.User, error) {
	q := &users.Query{
		Email: email,
	}

	return us.Store.UserGet(q)
}

// UserCreate ...
func (us *Users) UserCreate(u *users.User) error {
	return us.Store.UserCreate(u)
}

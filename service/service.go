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

	return us.Store.Get(q)
}

// GetByEmail ...
func (us *Users) GetByEmail(email string) (*users.User, error) {
	q := &users.Query{
		Email: email,
	}

	return us.Store.Get(q)
}

// Create ...
func (us *Users) Create(u *users.User) error {
	return us.Store.Create(u)
}

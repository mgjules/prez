package user

import (
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mgjules/prez/gonertia-demo/internal/errors"
)

type repository struct {
	users sync.Map
}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) Seed(n uint16) {
	for range n {
		user := User{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
		}

		_, _ = r.Create(user)
	}
}

func (r *repository) Create(user User) (*User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.users.Store(user.ID, user)

	return &user, nil
}

func (r *repository) List() ([]User, error) {
	var users []User
	r.users.Range(func(_, value any) bool {
		user, ok := value.(User)
		if !ok {
			return true
		}

		users = append(users, user)

		return true
	})

	return users, nil
}

func (r *repository) Find(id uuid.UUID) (*User, error) {
	raw, found := r.users.Load(id)
	if !found {
		return nil, errors.ErrNotFound
	}

	user, ok := raw.(User)
	if !ok {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}

func (r *repository) Update(user User) (*User, error) {
	if user.ID == uuid.Nil {
		return nil, errors.ErrNotFound
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	old, err := r.Find(user.ID)
	if err != nil {
		return nil, err
	}

	user.Email = old.Email
	user.CreatedAt = old.CreatedAt
	user.UpdatedAt = time.Now().UTC()

	r.users.Store(user.ID, user)

	return &user, nil
}

func (r *repository) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.ErrNotFound
	}

	if _, found := r.users.Load(id); !found {
		return errors.ErrNotFound
	}

	r.users.Delete(id)

	return nil
}

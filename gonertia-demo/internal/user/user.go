package user

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/mgjules/prez/gonertia-demo/internal/validate"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) Validate() error {
	var verr validate.Error
	if u.Name == "" {
		verr.Add("name", "Name is required")
	}
	if u.Email == "" {
		verr.Add("email", "Email is required")
	} else if _, err := mail.ParseAddress(u.Email); err != nil {
		verr.Add("email", fmt.Sprintf("Email is invalid: %s", err.Error()))
	}
	if verr.HasErrors() {
		return &verr
	}

	return nil
}

package user_entity

import (
	"context"
	"l03/internal/internal_error"

	"github.com/google/uuid"
)

type User struct {
	ID   string
	Name string
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user User) *internal_error.InternalError
	FindById(ctx context.Context, id string) (*User, *internal_error.InternalError)
	FindUsers(ctx context.Context) ([]User, *internal_error.InternalError)
}

func Create(name string) (*User, *internal_error.InternalError) {
	user := &User{
		ID:   uuid.New().String(),
		Name: name,
		//Timestamp:   time.Now(),
	}

	if err := user.Validade(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Validade() *internal_error.InternalError {
	if len(u.Name) < 5 {
		return internal_error.NewBadRequestError("Invalid name")
	}

	return nil
}

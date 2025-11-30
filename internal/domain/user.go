package domain

import "context"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type UserUsecase interface {
	Register(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
}

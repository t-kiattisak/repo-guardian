package usecase

import (
	"context"
	"time"

	"repo-guardian/internal/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(u domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (a *userUsecase) Register(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.userRepo.Create(ctx, user)
}

func (a *userUsecase) GetUser(c context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.userRepo.GetByID(ctx, id)
}

func (a *userUsecase) DeleteUser(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.userRepo.Delete(ctx, id)
}

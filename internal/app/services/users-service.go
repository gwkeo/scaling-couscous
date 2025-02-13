package services

import (
	"context"
	"errors"
	appErrors "github.com/gwkeo/scaling-couscous/internal/app/errors"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
)

type Repo interface {
	CreateUser(context.Context, *models.User) (int64, error)
	User(context.Context, int64) (*models.User, error)
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, int64) error
	UserByEmail(context.Context, string) (*models.User, error)
}

type UserService struct {
	repo Repo
}

func NewUsersService(repo Repo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *models.User) (int64, error) {
	_, err := s.repo.UserByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, appErrors.ErrUserNotFound) {
			return 0, err
		}
	}

	res, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Delete(ctx context.Context, user int64) error {
	err := s.repo.DeleteUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ReadByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Read(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.User(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

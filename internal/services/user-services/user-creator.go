package user_services

import (
	"context"
	"os/user"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *user.User) (int64, error)
}

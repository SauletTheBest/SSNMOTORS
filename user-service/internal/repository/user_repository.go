package repository

import (
    "context"
    "user-service/internal/model"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) (string, error)
    FindByID(ctx context.Context, id string) (*model.User, error)
    FindByUsername(ctx context.Context, username string) (*model.User, error)
}
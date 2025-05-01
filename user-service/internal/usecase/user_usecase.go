// internal/usecase/user_usecase.go
package usecase

import (
    "context"
    "errors"
    "user-service/internal/model"
    "user-service/internal/repository"

    "golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
    repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
    return &UserUsecase{repo: repo}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *model.User) (string, error) {
    if user.Username == "" || user.Email == "" {
        return "", errors.New("username and email are required")
    }
    // hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return "", errors.New("failed to hash password")
    }
    user.Password = string(hash)
    return u.repo.Create(ctx, user)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*model.User, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }
    return u.repo.FindByID(ctx, id)
}

func (u *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    if username == "" {
        return nil, errors.New("username is required")
    }
    return u.repo.FindByUsername(ctx, username)
}

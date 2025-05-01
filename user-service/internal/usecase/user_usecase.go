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

func (u *UserUsecase) Register(ctx context.Context, email, password, name string) (*model.User, error) {
    if email == "" || password == "" || name == "" {
        return nil, errors.New("missing fields")
    }
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    user := &model.User{
        Email:    email,
        Password: string(hashed),
        Name:     name,
        Role:     "user",
    }
    id, err := u.repo.Create(ctx, user)
    if err != nil {
        return nil, err
    }
    user.ID = id
    return user, nil
}

func (u *UserUsecase) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
    user, err := u.repo.FindByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    return user, nil
}

func (u *UserUsecase) GetProfile(ctx context.Context, id string) (*model.User, error) {
    return u.repo.FindByID(ctx, id)
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, id, name, email string) (*model.User, error) {
    user, err := u.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    user.Name = name
    user.Email = email
    if err := u.repo.Update(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

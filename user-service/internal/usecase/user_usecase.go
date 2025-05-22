package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
    repo  repository.UserRepository
    cache *redis.Client
}

func NewUserUsecase(repo repository.UserRepository, cache *redis.Client) *UserUsecase {
    return &UserUsecase{repo: repo, cache: cache}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *model.User) (string, error) {
    if user.Username == "" || user.Email == "" {
        return "", errors.New("username and email are required")
    }

    session, err := u.repo.StartSession(ctx)
    if err != nil {
        return "", err
    }
    defer session.EndSession(ctx)

    result, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }
    user.Password = string(hash)

    userID, err := u.repo.CreateWithTx(sessCtx, sessCtx, user) // <== оба ctx
    if err != nil {
        return nil, err
    }

    err = u.repo.LogAction(sessCtx, sessCtx, "create_user", userID)
    if err != nil {
        return nil, err
    }

    return userID, nil
	})

    if err != nil {
        return "", err
    }

    return result.(string), nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*model.User, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }

    cacheKey := "user_profile:" + id
    cachedUser, err := u.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var user model.User
        if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
            log.Println("User from cache")
            return &user, nil
        }
    }

    user, err := u.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    data, _ := json.Marshal(user)
    u.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    log.Println("User from mongo")
    return user, nil
}

func (u *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    if username == "" {
        return nil, errors.New("username is required")
    }
    return u.repo.FindByUsername(ctx, username)
}

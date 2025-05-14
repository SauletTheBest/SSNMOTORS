// internal/usecase/user_usecase.go
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
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo  repository.UserRepository
	cache *redis.Client //cache client
}

func NewUserUsecase(repo repository.UserRepository, cache *redis.Client) *UserUsecase {
	return &UserUsecase{repo: repo, cache: cache}
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

	// Check Redis cache first
	cacheKey := "user_profile:" + id
	cachedUser, err := u.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var user model.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			// If found in cache, return it
			log.Println("User from cache")
			return &user, nil
		}
	}

	// If not found in cache, query MongoDB
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		log.Println("User from mongo")
		return nil, err
	}

	// Cache the result in Redis
	data, _ := json.Marshal(user)
	u.cache.Set(ctx, cacheKey, data, 5*time.Minute) // Set cache expiration for 5 minutes

	return user, nil
}

func (u *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	return u.repo.FindByUsername(ctx, username)
}

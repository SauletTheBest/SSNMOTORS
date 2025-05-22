package repository

import (
	"context"
	"user-service/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) (string, error)
    FindByID(ctx context.Context, id string) (*model.User, error)
    FindByUsername(ctx context.Context, username string) (*model.User, error)
    StartSession(ctx context.Context) (mongo.Session, error)
    CreateWithTx(ctx context.Context, sessionCtx mongo.SessionContext, user *model.User) (string, error)
    LogAction(ctx context.Context, sessionCtx mongo.SessionContext, action string, userID string) error
}
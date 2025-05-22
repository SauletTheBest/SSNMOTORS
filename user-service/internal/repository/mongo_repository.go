// internal/repository/mongo_user_repository.go
package repository

import (
	"context"
	"errors"
	"user-service/internal/model"
    "time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
    coll *mongo.Collection
}

func NewMongoUserRepository(coll *mongo.Collection) *MongoUserRepository {
    // ensure unique index on username
    coll.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    )
    return &MongoUserRepository{coll: coll}
}

func (r *MongoUserRepository) StartSession(ctx context.Context) (mongo.Session, error) {
    return r.coll.Database().Client().StartSession()
}

func (r *MongoUserRepository) Create(ctx context.Context, user *model.User) (string, error) {
    obj := bson.M{"username": user.Username, "password": user.Password, "email": user.Email}
    res, err := r.coll.InsertOne(ctx, obj)
    if err != nil {
        if we, ok := err.(mongo.WriteException); ok {
            for _, e := range we.WriteErrors {
                if e.Code == 11000 {
                    return "", errors.New("username already exists")
                }
            }
        }
        return "", err
    }
    id := res.InsertedID.(primitive.ObjectID)
    return id.Hex(), nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid id format")
    }
    var u struct {
        ID       primitive.ObjectID `bson:"_id"`
        Username string             `bson:"username"`
        Password string             `bson:"password"`
        Email    string             `bson:"email"`
    }
    err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return &model.User{ID: u.ID.Hex(), Username: u.Username, Password: u.Password, Email: u.Email}, nil
}

func (r *MongoUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    var u struct {
        ID       primitive.ObjectID `bson:"_id"`
        Username string             `bson:"username"`
        Password string             `bson:"password"`
        Email    string             `bson:"email"`
    }
    err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return &model.User{ID: u.ID.Hex(), Username: u.Username, Password: u.Password, Email: u.Email}, nil
}

func (r *MongoUserRepository) LogAction(ctx context.Context, sessionCtx mongo.SessionContext, action string, userID string) error {
    logEntry := bson.M{
        "action":    action,
        "user_id":   userID,
        "timestamp": time.Now(),
    }
    _, err := r.coll.Database().Collection("logs").InsertOne(sessionCtx, logEntry)
    return err
}

func (r *MongoUserRepository) CreateWithTx(ctx context.Context, sessionCtx mongo.SessionContext, user *model.User) (string, error) {
    obj := bson.M{"username": user.Username, "password": user.Password, "email": user.Email}
    res, err := r.coll.InsertOne(sessionCtx, obj)
    if err != nil {
        return "", err
    }
    id := res.InsertedID.(primitive.ObjectID)
    return id.Hex(), nil
}
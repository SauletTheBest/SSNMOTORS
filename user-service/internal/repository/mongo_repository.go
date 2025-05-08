package repository

import (
    "context"
    "errors"
    "user-service/internal/model"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
    coll *mongo.Collection
}

func NewMongoUserRepository(coll *mongo.Collection) *MongoUserRepository {
    coll.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    )
    return &MongoUserRepository{coll: coll}
}

// DTO для работы с Mongo
type userDTO struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Email    string             `bson:"email"`
    Password string             `bson:"password"`
    Name     string             `bson:"name"`
    Role     string             `bson:"role"`
}

func (r *MongoUserRepository) Create(ctx context.Context, user *model.User) (string, error) {
    dto := userDTO{
        Email:    user.Email,
        Password: user.Password,
        Name:     user.Name,
        Role:     user.Role,
    }
    res, err := r.coll.InsertOne(ctx, dto)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) {
            return "", errors.New("email already exists")
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
    var dto userDTO
    err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&dto)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &model.User{
        ID:       dto.ID.Hex(),
        Email:    dto.Email,
        Password: dto.Password,
        Name:     dto.Name,
        Role:     dto.Role,
    }, nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var dto userDTO
    err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&dto)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &model.User{
        ID:       dto.ID.Hex(),
        Email:    dto.Email,
        Password: dto.Password,
        Name:     dto.Name,
        Role:     dto.Role,
    }, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *model.User) error {
    oid, err := primitive.ObjectIDFromHex(user.ID)
    if err != nil {
        return errors.New("invalid id format")
    }

    update := bson.M{
        "$set": bson.M{
            "name":  user.Name,
            "email": user.Email,
        },
    }
    _, err = r.coll.UpdateByID(ctx, oid, update)
    return err
}

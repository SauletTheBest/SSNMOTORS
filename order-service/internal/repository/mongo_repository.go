package repository

import (
	"context"
	"errors"
	"order-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOrderRepository struct {
	coll *mongo.Collection
}

func NewMongoOrderRepository(coll *mongo.Collection) *MongoOrderRepository {
	return &MongoOrderRepository{coll: coll}
}

type orderDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	CarID     string             `bson:"car_id"`
	Quantity  int32              `bson:"quantity"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (r *MongoOrderRepository) Create(ctx context.Context, order *model.Order) (string, error) {
	dto := orderDTO{
		UserID:    order.UserID,
		CarID:     order.CarID,
		Quantity:  order.Quantity,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}
	res, err := r.coll.InsertOne(ctx, dto)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (r *MongoOrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid order id")
	}
	var dto orderDTO
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &model.Order{
		ID:        dto.ID.Hex(),
		UserID:    dto.UserID,
		CarID:     dto.CarID,
		Quantity:  dto.Quantity,
		Status:    dto.Status,
		CreatedAt: dto.CreatedAt,
	}, nil
}

func (r *MongoOrderRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*model.Order
	for cursor.Next(ctx) {
		var dto orderDTO
		if err := cursor.Decode(&dto); err != nil {
			continue
		}
		orders = append(orders, &model.Order{
			ID:        dto.ID.Hex(),
			UserID:    dto.UserID,
			CarID:     dto.CarID,
			Quantity:  dto.Quantity,
			Status:    dto.Status,
			CreatedAt: dto.CreatedAt,
		})
	}
	return orders, nil
}

func (r *MongoOrderRepository) Update(ctx context.Context, order *model.Order) error {
	oid, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return errors.New("invalid order id")
	}

	update := bson.M{
		"$set": bson.M{
			"status": order.Status,
		},
	}

	_, err = r.coll.UpdateByID(ctx, oid, update)
	return err
}

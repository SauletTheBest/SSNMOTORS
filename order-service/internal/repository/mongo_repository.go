package repository

import (
	"context"
	"errors"
	"order-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOrderRepository struct {
	collection *mongo.Collection
}

func NewMongoOrderRepository(collection *mongo.Collection) *MongoOrderRepository {
	return &MongoOrderRepository{collection: collection}
}

func (r *MongoOrderRepository) Create(ctx context.Context, order *model.Order) (string, error) {
	doc := bson.M{
		"user_id": order.UserID,
		"products": bson.A{},
		"total":   order.Total,
		"status":  order.Status,
	}
	for _, p := range order.Products {
		doc["products"] = append(doc["products"].(bson.A), bson.M{
			"product_id": p.ProductID,
			"quantity":   p.Quantity,
		})
	}
	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (r *MongoOrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID")
	}
	var result struct {
		ID        primitive.ObjectID `bson:"_id"`
		UserID    string             `bson:"user_id"`
		Products  []struct {
			ProductID string `bson:"product_id"`
			Quantity  int    `bson:"quantity"`
		} `bson:"products"`
		Total  float64 `bson:"total"`
		Status string  `bson:"status"`
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	order := &model.Order{
		ID:     result.ID.Hex(),
		UserID: result.UserID,
		Total:  result.Total,
		Status: result.Status,
	}
	for _, p := range result.Products {
		order.Products = append(order.Products, model.Product{
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
		})
	}
	return order, nil
}

func (r *MongoOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID")
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *MongoOrderRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var orders []*model.Order
	for cursor.Next(ctx) {
		var result struct {
			ID        primitive.ObjectID `bson:"_id"`
			UserID    string             `bson:"user_id"`
			Products  []struct {
				ProductID string `bson:"product_id"`
				Quantity  int    `bson:"quantity"`
			} `bson:"products"`
			Total  float64 `bson:"total"`
			Status string  `bson:"status"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		order := &model.Order{
			ID:     result.ID.Hex(),
			UserID: result.UserID,
			Total:  result.Total,
			Status: result.Status,
		}
		for _, p := range result.Products {
			order.Products = append(order.Products, model.Product{
				ProductID: p.ProductID,
				Quantity:  p.Quantity,
			})
		}
		orders = append(orders, order)
	}
	return orders, nil
}
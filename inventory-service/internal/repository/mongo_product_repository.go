package repository

import (
    "context"
    "errors"
    "inventory-service/internal/model"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProductRepository struct {
    coll *mongo.Collection
}

func NewMongoProductRepository(coll *mongo.Collection) *MongoProductRepository {
    // Убедимся, что есть уникальный индекс по name (опционально)
    coll.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
            Keys:    bson.D{{Key: "name", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    )
    return &MongoProductRepository{coll: coll}
}

func (r *MongoProductRepository) Create(ctx context.Context, p *model.Product) (string, error) {
    doc := bson.M{
        "name":        p.Name,
        "description": p.Description,
        "category":    p.Category,
        "stock":       p.Stock,
        "price":       p.Price,
    }
    res, err := r.coll.InsertOne(ctx, doc)
    if err != nil {
        if we, ok := err.(mongo.WriteException); ok {
            for _, e := range we.WriteErrors {
                if e.Code == 11000 {
                    return "", errors.New("product already exists")
                }
            }
        }
        return "", err
    }
    oid := res.InsertedID.(primitive.ObjectID)
    return oid.Hex(), nil
}

func (r *MongoProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid ID format")
    }
    var doc bson.M
    err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("product not found")
    } else if err != nil {
        return nil, err
    }
    return &model.Product{
        ID:          id,
        Name:        doc["name"].(string),
        Description: doc["description"].(string),
        Category:    doc["category"].(string),
        Stock:       int32(doc["stock"].(int32)),
        Price:       doc["price"].(float64),
    }, nil
}

func (r *MongoProductRepository) Update(ctx context.Context, p *model.Product) error {
    oid, err := primitive.ObjectIDFromHex(p.ID)
    if err != nil {
        return errors.New("invalid ID format")
    }
    update := bson.M{
        "name":        p.Name,
        "description": p.Description,
        "category":    p.Category,
        "stock":       p.Stock,
        "price":       p.Price,
    }
    res, err := r.coll.UpdateByID(ctx, oid, bson.M{"$set": update})
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return errors.New("product not found")
    }
    return nil
}

func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid ID format")
    }
    res, err := r.coll.DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return errors.New("product not found")
    }
    return nil
}

func (r *MongoProductRepository) List(ctx context.Context, category string, page, limit int32) ([]*model.Product, error) {
    filter := bson.M{}
    if category != "" {
        filter["category"] = category
    }

    opts := options.Find().
        SetSkip(int64((page-1)*limit)).
        SetLimit(int64(limit))

    // Структура, точно соответствующая полям в БД
    type dbProduct struct {
        ID          primitive.ObjectID `bson:"_id"`
        Name        string             `bson:"name"`
        Description string             `bson:"description"`
        Category    string             `bson:"category"`
        Stock       int32              `bson:"stock"`
        Price       float64            `bson:"price"`
    }

    cursor, err := r.coll.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var out []*model.Product
    for cursor.Next(ctx) {
        var dp dbProduct
        if err := cursor.Decode(&dp); err != nil {
            // пропускаем некорректный документ
            continue
        }
        out = append(out, &model.Product{
            ID:          dp.ID.Hex(),
            Name:        dp.Name,
            Description: dp.Description,
            Category:    dp.Category,
            Stock:       dp.Stock,
            Price:       dp.Price,
        })
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    return out, nil
}
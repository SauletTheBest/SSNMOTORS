package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// AutoInc manages auto-incremented IDs for any collection
type AutoInc struct {
	collection *mongo.Collection
}

// Counter document structure in MongoDB
type Counter struct {
	ID      string `bson:"_id"`     // Collection name as _id
	Counter uint64 `bson:"counter"` // Auto-incremented value
}

// NewAutoInc is Constructor of AutoInc
func NewAutoInc(db *mongo.Database) *AutoInc {
	return &AutoInc{
		collection: db.Collection(CollectionAutoInc),
	}
}

// Next -> next auto-incremented ID for the given collection
func (a *AutoInc) Next(ctx context.Context, coll string) (uint64, error) {
	log.Printf("Generating next ID for collection: %s\n", coll)

	// Define the filter and update operations
	filter := bson.M{"_id": coll}
	update := bson.M{"$inc": bson.M{"counter": 1}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).                 // Create the document if it doesn’t exist
		SetReturnDocument(options.After) // Return the updated document

	// Perform the atomic update and retrieve the result
	var result Counter
	err := a.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This case shouldn’t happen with Upsert(true), but included for safety
			return 1, nil
		}
		return 0, fmt.Errorf("FindOneAndUpdate failed: %w", err)
	}

	return result.Counter, nil
}

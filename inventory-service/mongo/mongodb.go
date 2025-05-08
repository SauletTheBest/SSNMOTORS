package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds the MongoDB connection and client.
type DB struct {
	Conn   *mongo.Database
	Client *mongo.Client
}

// NewDB creates a connection to MongoDB and returns a DB struct.
func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	// Use the existing genConnectURL method to build the connection string
	connectURL := cfg.genConnectURL()

	// Set up client options
	clientOptions := options.Client().ApplyURI(connectURL)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Select the database
	db := &DB{
		Conn:   client.Database(cfg.Database),
		Client: client,
	}

	// Verify the connection with a ping
	err = db.Client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Start a background reconnection routine
	go db.reconnectOnFailure(ctx, cfg)

	log.Println("Connected to MongoDB successfully")
	return db, nil
}

// reconnectOnFailure attempts to reconnect if the connection is lost.
func (db *DB) reconnectOnFailure(ctx context.Context, cfg Config) {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			err := db.Client.Ping(ctx, nil)
			if err != nil {
				log.Printf("Lost connection to MongoDB: %v", err)
				// Attempt to reconnect using the same config
				newClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.genConnectURL()))
				if err == nil {
					db.Client = newClient
					db.Conn = newClient.Database(cfg.Database)
					log.Println("Reconnected to MongoDB successfully")
				} else {
					log.Printf("Failed to reconnect to MongoDB: %v", err)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			err := db.Client.Disconnect(ctx)
			if err != nil {
				log.Printf("Failed to close MongoDB connection: %v", err)
			} else {
				log.Println("MongoDB connection closed successfully")
			}
			return
		}
	}
}

// Ping checks the connection to MongoDB.
func (db *DB) Ping(ctx context.Context) error {
	err := db.Client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("MongoDB connection error: %w", err)
	}
	return nil
}

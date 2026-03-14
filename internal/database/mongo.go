package database

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type DBEngine struct {
	Client *mongo.Client 
}

func Connect(uri string) *DBEngine {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	clientOps := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOps)
	if err != nil {
		slog.Error("Failed to create MongoDB client",
			slog.String("error", err.Error()),
		)
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error("Failed to ping mongodb", slog.String("error", err.Error()))
		panic(err)
	}

	slog.Info("Successfully connected to MongoDB!")

	return &DBEngine{
		Client: client,
	}
}

func (db *DBEngine) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}

func (db *DBEngine) GetCollection(collectionName string) *mongo.Collection {
	// HARDCODING THE DB NAME TO "restaurant"
	// TODO: PULL THIS FROM ENV FILE
	return db.Client.Database("restaurant").Collection(collectionName)
}
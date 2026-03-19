package database

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

)

type Database struct {
	Client *mongo.Client
	DBName string
}

func Connect(uri, dbName string) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	return &Database{
		Client: client,
		DBName: dbName,
	}
}

func (db *Database) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}

func (db *Database) GetCollection(collectionName string) *mongo.Collection {
	return db.Client.Database(db.DBName).Collection(collectionName)
}

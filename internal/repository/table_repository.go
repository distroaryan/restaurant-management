package repository

import (
	"context"
	"errors"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TableRepository struct {
	collection *mongo.Collection
}

func NewTableRepositroy(db *database.Database) *TableRepository {
	return &TableRepository{
		collection: db.GetCollection("tables"),
	}
}

func (r *TableRepository) CreateTable(ctx context.Context, table *models.Table) error {
	table.Status = models.TableStatusAvailable

	resp, err := r.collection.InsertOne(ctx, &table)
	if err != nil {
		return err
	}
	if objectId, ok := resp.InsertedID.(bson.ObjectID); ok {
		table.ID = objectId
	}
	return nil
}

func (r *TableRepository) GetTableById(ctx context.Context, id string) (*models.Table, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectId}
	var table models.Table

	err = r.collection.FindOne(ctx, filter).Decode(&table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func (r *TableRepository) GetAllTables(ctx context.Context) ([]*models.Table, error) {
	filter := bson.M{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tables []*models.Table
	if err := cursor.All(ctx, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

// BookTable atomically updates the table status to full and assigns it to a user
func (r *TableRepository) BookTable(ctx context.Context, tableId string, userId string) error {
	objectId, err := bson.ObjectIDFromHex(tableId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectId,
		"status": models.TableStatusAvailable,
	}

	update := bson.M{
		"$set": bson.M{
			"status": models.TableStatusFull,
			"user_id": userId,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("table not found or already booked")
	}
	return nil
}

// ReleaseTable frees a table so it can be booked again
func (r *TableRepository) ReleaseTable(ctx context.Context, tableId string) error {
	objectId, err := bson.ObjectIDFromHex(tableId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectId,
	}

	update := bson.M{
		"$set": bson.M{
			"status": models.TableStatusAvailable,
		},
		"$unset": bson.M{
			"user_id": "",
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("table not found")
	}
	return nil
}

func (r *TableRepository) UpdateTableStatus(ctx context.Context, tableId string, tableStatus models.TableStatus) error {
	objectId, err := bson.ObjectIDFromHex(tableId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectId,
	}

	update := bson.M{
		"$set": bson.M{"status": tableStatus},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

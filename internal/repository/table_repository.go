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

func NewTableRepositroy(db *database.DBEngine) *TableRepository {
	return &TableRepository{
		collection: db.GetCollection("tables"),
	}
} 

func (r *TableRepository) CreateTable(ctx context.Context, table *models.Table) error {
	table.ReservedSeats = 0
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

// BookSeats atomically increments the reserved seats for a given table,
// avoiding race conditions by validating that there is enough capacity

func (r *TableRepository) BookSeats(ctx context.Context, tableId string, seats int) error {
	objectId, err := bson.ObjectIDFromHex(tableId)
	if err != nil {
		return err 
	}
	
	// make sure seats are available
	// capacity - reservedSeats >= seats
	filter := bson.M{
		"_id": objectId,
		"$expr": bson.M{
			"$gte": []interface{}{
				bson.M{"$subtract": []string{"$capacity", "$reserved_seats"}},
				seats, 
			},
		},
	}

	// update query
	update := bson.M{
		"$inc": bson.M{"reserved_seats": seats},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err 
	}
	if res.MatchedCount == 0 {
		return errors.New("table not found or insufficient capacity available")
	}
	return nil 
}

func (r *TableRepository) ReleaseSeats(ctx context.Context, tableId string, seatsToRelease int) error {
	objectId, err := bson.ObjectIDFromHex(tableId)
	if err != nil {
		return err 
	}

	// check reserved_seats >= seatsToRelease
	filter := bson.M{
		"_id": objectId,
		"reserved_seats": bson.M{"$gte": seatsToRelease},
	}

	update := bson.M{
		"$inc": bson.M{"reserved_seats": -seatsToRelease},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err 
	}
	if res.MatchedCount == 0 {
		return errors.New("table not found or invalid seatsToRealease sent")
	}
	return nil 
}

func (r *TableRepository) UpdateTableStatus(ctx context.Context, tableId string, tableStatus models.TableStatus) error {
	objectId ,err := bson.ObjectIDFromHex(tableId)
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
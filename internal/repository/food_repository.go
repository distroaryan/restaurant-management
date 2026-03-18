package repository

import (
	"context"

	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FoodRepository struct {
	collection *mongo.Collection
}

func NewFoodRepository(db *database.Database) *FoodRepository {
	return &FoodRepository{
		collection: db.GetCollection("foods"),
	}
}

func (r *FoodRepository) CreateFood(ctx context.Context, food *models.Food) error {
	resp, err := r.collection.InsertOne(ctx, food)
	if err != nil {
		return err
	}
	if oid, ok := resp.InsertedID.(bson.ObjectID); ok {
		food.ID = oid
	}
	return nil
}

func (r *FoodRepository) GetFoodById(ctx context.Context, id string) (*models.Food, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectId}
	var food models.Food
	err = r.collection.FindOne(ctx, filter).Decode(&food)
	if err != nil {
		return nil, err
	}
	return &food, err
}

func (r *FoodRepository) GetFoodByMenu(ctx context.Context, menuID string) ([]*models.Food, error) {
	objectId, err := bson.ObjectIDFromHex(menuID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"menu_id": objectId}
	var foods []*models.Food
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &foods); err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *FoodRepository) GetAllFoods(ctx context.Context) ([]*models.Food, error) {
	var foods []*models.Food
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &foods); err != nil {
		return nil, err
	}
	return foods, nil
}

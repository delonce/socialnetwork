package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/delonce/socialnetwork/internal/database"
	"github.com/delonce/socialnetwork/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	colls  map[string]*mongo.Collection
	logger *logging.Logger
}

func NewStorage(collectionMap map[string]*mongo.Collection, logger *logging.Logger) database.DBStorage {
	return &mongoDB{
		colls:  collectionMap,
		logger: logger,
	}
}

func (db *mongoDB) CreateObject(ctx context.Context, model interface{}, key string) (string, error) {
	result, err := db.colls[key].InsertOne(ctx, model)

	if err != nil {
		db.logger.Errorf("Creating object ends with error: %v", err)
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	db.logger.Errorf("Failed to create ObjectId for %s", model)

	return "", fmt.Errorf("Failed to create object %s", model)
}

func (db *mongoDB) FindOneObject(ctx context.Context, filter bson.M, key string) (*mongo.SingleResult, error) {
	result := db.colls[key].FindOne(ctx, filter)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("Object not found")
		}

		db.logger.Errorf("Failed to find object with error: %v", result.Err())
		return nil, result.Err()
	}

	return result, nil
}

func (db *mongoDB) FindObjects(ctx context.Context, filter bson.M, key string, findOpt ...*options.FindOptions) (*mongo.Cursor, error) {
	result, err := db.colls[key].Find(ctx, filter, findOpt...)

	if err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("Object not found")
		}

		db.logger.Errorf("Failed to find object with error: %v", result.Err())
		return nil, result.Err()
	}

	return result, nil
}

func (db *mongoDB) CountObjects(ctx context.Context, filter bson.M, key string) (int64, error) {
	result, err := db.colls[key].CountDocuments(ctx, filter)

	if err != nil {
		db.logger.Errorf("Failed to find object with error: %v", err)
		return 0, err
	}

	return result, nil
}

func (db *mongoDB) Update(ctx context.Context, filter bson.M, model interface{}, key string) error {
	update := bson.D{{Key: "$set", Value: model}}
	result, err := db.colls[key].UpdateOne(ctx, filter, update)

	if err != nil {
		db.logger.Errorf("Failed to execute update user query, error: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("Object not found")
	}

	db.logger.Infof("Model %s with filter %s updated", model, filter)
	return nil
}

func (db *mongoDB) Delete(ctx context.Context, filter bson.M, key string) error {
	result, err := db.colls[key].DeleteOne(ctx, filter)

	if err != nil {
		db.logger.Errorf("Failed to delete %s with error: %v", filter, err)
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("Object not found")
	}

	return nil
}

func (db *mongoDB) Aggregate(ctx context.Context, pipeline interface{}, key string) (*mongo.Cursor, error) {
	result, err := db.colls[key].Aggregate(ctx, pipeline)

	if err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("Object not found")
		}

		db.logger.Errorf("Failed to find object with error: %v", result.Err())
		return nil, result.Err()
	}

	return result, nil
}

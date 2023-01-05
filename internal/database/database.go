package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBStorage interface {
	CreateObject(ctx context.Context, model interface{}, key string) (string, error)
	FindObjects(ctx context.Context, filter bson.M, key string, findOpt ...*options.FindOptions) (*mongo.Cursor, error)
	FindOneObject(ctx context.Context, filter bson.M, key string) (*mongo.SingleResult, error)
	Update(ctx context.Context, filter bson.M, model interface{}, key string) error
	Delete(ctx context.Context, filter bson.M, key string) error
	CountObjects(ctx context.Context, filter bson.M, key string) (int64, error)
	Aggregate(ctx context.Context, pipeline interface{}, key string) (*mongo.Cursor, error)
}

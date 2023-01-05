package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBClient(ctx context.Context, host, port, username, password, database, authDB string) (*mongo.Database, error) {
	var mongoAddr string

	if username == "" && password == "" {
		mongoAddr = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		mongoAddr = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	clientOpts := options.Client().ApplyURI(mongoAddr)
	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db with error: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to connect to db with error: %v", err)
	}

	return client.Database(database), nil
}

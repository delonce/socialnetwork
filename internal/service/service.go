package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/pkg/dbclient/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_COLLECTION           = "users"
	SESSION_COLLECTION        = "sessions"
	FRIEND_REQUEST_COLLECTION = "friend_requests"
	MESSAGE_COLLECTION        = "messages"
)

func InitNewDatabase(config *config.Config) (*mongo.Database, error) {
	dbConfig := config.Database

	database, err := mongodb.NewDBClient(context.Background(), dbConfig.Addr, strconv.Itoa(int(dbConfig.Port)),
		dbConfig.Login, dbConfig.Password, dbConfig.DBName, dbConfig.DBName)

	if err != nil {
		return nil, fmt.Errorf("Error occcurs when creating database client. %v", err)
	}

	return database, nil
}

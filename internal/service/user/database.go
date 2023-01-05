package user

import (
	"context"
	"time"

	"github.com/delonce/socialnetwork/internal/database"
	"github.com/delonce/socialnetwork/internal/database/mongodb"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserQueries interface {
	CreateNewUser(ctx context.Context, user *service.User) (string, error)
	FindUserByUsername(ctx context.Context, username string) (*mongo.SingleResult, error)
	FindUserByEmail(ctx context.Context, email string) (*mongo.SingleResult, error)
	FindUserByID(ctx context.Context, userID string) (*mongo.SingleResult, error)
	FindUserByCredentials(ctx context.Context, username, passwordHash string) (*mongo.SingleResult, error)
	DeleteUser(ctx context.Context, user *service.User) error

	AddRefreshToken(ctx context.Context, refreshToken *service.RefreshToken) (string, error)
	FindRefreshTokenByUUID(ctx context.Context, refreshTokenUUID string) (*mongo.SingleResult, error)
	FindExpiredRefreshTokens(ctx context.Context, userID string) (*mongo.Cursor, error)
	DeleteRefreshToken(ctx context.Context, refreshTokenUUID string) error
}

type UserDB struct {
	Storage database.DBStorage
	Logger  *logging.Logger
}

func NewUserDB(logger *logging.Logger, database *mongo.Database) UserQueries {
	storage := mongodb.NewStorage(
		map[string]*mongo.Collection{
			service.USER_COLLECTION:    database.Collection(service.USER_COLLECTION),
			service.SESSION_COLLECTION: database.Collection(service.SESSION_COLLECTION),
		},
		logger,
	)

	return &UserDB{
		Storage: storage,
		Logger:  logger,
	}
}

func (userStorage *UserDB) CreateNewUser(ctx context.Context, user *service.User) (string, error) {
	st := userStorage.Storage

	return st.CreateObject(ctx, user, service.USER_COLLECTION)
}

func (userStorage *UserDB) AddRefreshToken(ctx context.Context, refreshToken *service.RefreshToken) (string, error) {
	st := userStorage.Storage

	return st.CreateObject(ctx, refreshToken, service.SESSION_COLLECTION)
}

func (userStorage *UserDB) FindUserByUsername(ctx context.Context, username string) (*mongo.SingleResult, error) {
	st := userStorage.Storage
	query := bson.M{"username": username}

	return st.FindOneObject(ctx, query, service.USER_COLLECTION)
}

func (userStorage *UserDB) FindUserByEmail(ctx context.Context, email string) (*mongo.SingleResult, error) {
	st := userStorage.Storage
	query := bson.M{"email": email}

	return st.FindOneObject(ctx, query, service.USER_COLLECTION)
}

func (userStorage *UserDB) FindUserByID(ctx context.Context, userID string) (*mongo.SingleResult, error) {
	st := userStorage.Storage
	objUserID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objUserID}

	return st.FindOneObject(ctx, query, service.USER_COLLECTION)
}

func (userStorage *UserDB) FindUserByCredentials(ctx context.Context, username, passwordHash string) (*mongo.SingleResult, error) {
	st := userStorage.Storage
	query := bson.M{"username": username, "password": passwordHash}

	return st.FindOneObject(ctx, query, service.USER_COLLECTION)
}

func (userStorage *UserDB) DeleteUser(ctx context.Context, user *service.User) error {
	st := userStorage.Storage
	query := bson.M{"_id": user.ID}

	return st.Delete(ctx, query, service.USER_COLLECTION)
}

func (userStorage *UserDB) DeleteRefreshToken(ctx context.Context, refreshTokenUUID string) error {
	st := userStorage.Storage
	query := bson.M{"uuid": refreshTokenUUID}

	return st.Delete(ctx, query, service.SESSION_COLLECTION)
}

func (userStorage *UserDB) FindRefreshTokenByUUID(ctx context.Context, refreshTokenUUID string) (*mongo.SingleResult, error) {
	st := userStorage.Storage
	query := bson.M{"uuid": refreshTokenUUID}

	return st.FindOneObject(ctx, query, service.SESSION_COLLECTION)
}

func (userStorage *UserDB) FindExpiredRefreshTokens(ctx context.Context, userID string) (*mongo.Cursor, error) {
	st := userStorage.Storage
	query := bson.M{
		"$and": []bson.M{ // you can try this in []interface
			{"userID": userID},
			{"expiresAt": bson.M{"$lt": time.Now()}},
		},
	}

	return st.FindObjects(ctx, query, service.SESSION_COLLECTION)
}

package friends

import (
	"context"

	"github.com/delonce/socialnetwork/internal/database"
	"github.com/delonce/socialnetwork/internal/database/mongodb"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendQueries interface {
	GetAllProbFriends(ctx context.Context, friendUsernames []string) (*mongo.Cursor, error)
	GetFriends(ctx context.Context, username string) (*mongo.Cursor, error)

	FindUserRequests(ctx context.Context, username string) (*mongo.Cursor, error)
	GetAllFriendRequestTo(ctx context.Context, to string) (*mongo.Cursor, error)
	GetRequestByNames(ctx context.Context, from string, to string) (*mongo.SingleResult, error)

	AddFriendRequest(ctx context.Context, request *service.FriendRequest) (string, error)
	DeleteRequest(ctx context.Context, from string, to string) error

	IsExistRequest(ctx context.Context, from string, to string) bool
	IsExistFriend(ctx context.Context, from string, to string) bool

	SetAcceptedMark(ctx context.Context, requestId string) error
	SetDeniedMark(ctx context.Context, requestID string, from string, to string) error

	GetUsersExcept(ctx context.Context, except []string) (*mongo.Cursor, error)
}

type FriendDB struct {
	Storage database.DBStorage
	Logger  *logging.Logger
}

func NewFriendDB(logger *logging.Logger, database *mongo.Database) FriendQueries {
	storage := mongodb.NewStorage(
		map[string]*mongo.Collection{
			service.USER_COLLECTION:           database.Collection(service.USER_COLLECTION),
			service.FRIEND_REQUEST_COLLECTION: database.Collection(service.FRIEND_REQUEST_COLLECTION),
		},
		logger,
	)

	return &FriendDB{
		Storage: storage,
		Logger:  logger,
	}
}

func (friendStorage *FriendDB) GetAllProbFriends(ctx context.Context, friendUsernames []string) (*mongo.Cursor, error) {
	st := friendStorage.Storage
	query := bson.M{
		"username": bson.M{"$nin": friendUsernames},
	}

	return st.FindObjects(ctx, query, service.USER_COLLECTION)
}

func (friendStorage *FriendDB) GetFriends(ctx context.Context, username string) (*mongo.Cursor, error) {
	query := bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"from": username},
				{"to": username},
			}},

			{"isaccepted": true},
		},
	}

	return friendStorage.Storage.FindObjects(ctx, query, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) AddFriendRequest(ctx context.Context, request *service.FriendRequest) (string, error) {
	st := friendStorage.Storage

	isReqExist := friendStorage.IsExistRequest(ctx, request.From, request.To)

	if isReqExist {
		return "", nil
	}

	isReqExist = friendStorage.IsExistRequest(ctx, request.To, request.From)

	if isReqExist {
		return "", nil
	}

	return st.CreateObject(ctx, request, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) IsExistRequest(ctx context.Context, from string, to string) bool {
	st := friendStorage.Storage

	query := bson.M{
		"$and": []bson.M{
			{"from": from},
			{"to": to},
		},
	}

	_, err := st.FindOneObject(ctx, query, service.FRIEND_REQUEST_COLLECTION)

	if err != nil {
		return false
	}

	return true
}

func (friendStorage *FriendDB) IsExistFriend(ctx context.Context, from string, to string) bool {
	st := friendStorage.Storage

	query := bson.M{
		"$or": []bson.M{
			{"$and": []bson.M{
				{"from": from},
				{"to": to},
				{"isaccepted": true},
				{"isdenied": false},
			}},

			{"$and": []bson.M{
				{"from": to},
				{"to": from},
				{"isaccepted": true},
				{"isdenied": false},
			}},
		},
	}

	_, err := st.FindOneObject(ctx, query, service.FRIEND_REQUEST_COLLECTION)

	if err != nil {
		return false
	}

	return true
}

func (friendStorage *FriendDB) DeleteRequest(ctx context.Context, from string, to string) error {
	st := friendStorage.Storage
	query := bson.M{
		"$and": []bson.M{
			{"from": from},
			{"to": to},
			{"isaccepted": false},
		},
	}

	return st.Delete(ctx, query, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) FindUserRequests(ctx context.Context, username string) (*mongo.Cursor, error) {
	st := friendStorage.Storage
	query := bson.M{
		"$and": []bson.M{
			{"to": username},
			{"isaccepted": false},
			{"isdenied": false},
		},
	}

	return st.FindObjects(ctx, query, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) GetAllFriendRequestTo(ctx context.Context, to string) (*mongo.Cursor, error) {
	st := friendStorage.Storage
	query := bson.M{
		"$and": []bson.M{
			{"to": to},
			{"isaccepted": false},
			{"isdenied": false},
		},
	}

	return st.FindObjects(ctx, query, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) SetAcceptedMark(ctx context.Context, requestID string) error {
	model := bson.M{
		"isaccepted": true,
		"isdenied":   false,
	}

	return friendStorage.SetMark(ctx, requestID, model)
}

func (friendStorage *FriendDB) SetDeniedMark(ctx context.Context, requestID string, from string, to string) error {
	model := bson.M{
		"from":       from,
		"to":         to,
		"isaccepted": false,
		"isdenied":   true,
	}

	return friendStorage.SetMark(ctx, requestID, model)
}

func (friendStorage *FriendDB) GetRequestByNames(ctx context.Context, from string, to string) (*mongo.SingleResult, error) {
	st := friendStorage.Storage
	query := bson.M{
		"$and": []bson.M{
			{"from": from},
			{"to": to},
		},
	}

	return st.FindOneObject(ctx, query, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) SetMark(ctx context.Context, requestID string, model bson.M) error {
	st := friendStorage.Storage
	objRequestID, err := primitive.ObjectIDFromHex(requestID)

	if err != nil {
		return err
	}

	query := bson.M{
		"_id": objRequestID,
	}

	return st.Update(ctx, query, model, service.FRIEND_REQUEST_COLLECTION)
}

func (friendStorage *FriendDB) GetUsersExcept(ctx context.Context, except []string) (*mongo.Cursor, error) {
	st := friendStorage.Storage

	query := bson.M{
		"username": bson.M{"$nin": except},
	}

	return st.FindObjects(ctx, query, service.USER_COLLECTION)
}

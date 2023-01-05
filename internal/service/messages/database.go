package messages

import (
	"context"

	"github.com/delonce/socialnetwork/internal/database"
	"github.com/delonce/socialnetwork/internal/database/mongodb"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageQueries interface {
	AddNewMessage(ctx context.Context, message service.Message) error

	GetAllDialogMessages(ctx context.Context, from string, to string) (*mongo.Cursor, error)
	GetLastMessage(ctx context.Context, username string, friend string) (*mongo.Cursor, error)

	CountAllNewMessages(ctx context.Context, username string) int64
	CountNewMessagesBySender(ctx context.Context, from string, to string) int64

	SetCheckMark(ctx context.Context, from string, to string) error
}

type MessageDB struct {
	Storage database.DBStorage
	Logger  *logging.Logger
}

func NewMessageDB(logger *logging.Logger, database *mongo.Database) MessageQueries {
	storage := mongodb.NewStorage(
		map[string]*mongo.Collection{
			service.MESSAGE_COLLECTION: database.Collection(service.MESSAGE_COLLECTION),
		},
		logger,
	)

	return &MessageDB{
		Storage: storage,
		Logger:  logger,
	}
}

func (msgDatabase *MessageDB) AddNewMessage(ctx context.Context, message service.Message) error {
	st := msgDatabase.Storage
	_, err := st.CreateObject(ctx, message, service.MESSAGE_COLLECTION)

	return err
}

func (msgDatabase *MessageDB) GetAllDialogMessages(ctx context.Context, from string, to string) (*mongo.Cursor, error) {
	st := msgDatabase.Storage

	query := bson.M{
		"$or": []bson.M{
			{"$and": []bson.M{
				{"from": from},
				{"to": to},
			}},

			{"$and": []bson.M{
				{"from": to},
				{"to": from},
			}},
		},
	}

	findOpts := options.FindOptions{}
	findOpts.SetSort(bson.D{{Key: "date", Value: 1}})

	return st.FindObjects(ctx, query, service.MESSAGE_COLLECTION, &findOpts)
}

func (msgDatabase *MessageDB) CountAllNewMessages(ctx context.Context, username string) int64 {
	st := msgDatabase.Storage

	query := bson.M{
		"$and": []bson.M{
			{"to": username},
			{"ischecked": false},
		},
	}

	num, err := st.CountObjects(ctx, query, service.MESSAGE_COLLECTION)

	if err != nil {
		msgDatabase.Logger.Panic("Error when counting, %v", err)
	}

	return num
}

func (msgDatabase *MessageDB) CountNewMessagesBySender(ctx context.Context, from string, to string) int64 {
	st := msgDatabase.Storage

	query := bson.M{
		"$and": []bson.M{
			{"from": from},
			{"to": to},
			{"ischecked": false},
		},
	}

	num, err := st.CountObjects(ctx, query, service.MESSAGE_COLLECTION)

	if err != nil {
		msgDatabase.Logger.Panic("Error when counting, %v", err)
	}

	return num
}

func (msgDatabase *MessageDB) GetLastMessage(ctx context.Context, username string, friend string) (*mongo.Cursor, error) {
	st := msgDatabase.Storage

	query := bson.M{
		"$or": []bson.M{
			{"$and": []bson.M{
				{"from": username},
				{"to": friend},
			}},

			{"$and": []bson.M{
				{"from": friend},
				{"to": username},
			}},
		},
	}

	findOpts := options.FindOptions{}

	findOpts.SetSort(bson.D{{Key: "date", Value: -1}})
	findOpts.SetLimit(1)

	return st.FindObjects(ctx, query, service.MESSAGE_COLLECTION, &findOpts)
}

func (msgDatabase *MessageDB) SetCheckMark(ctx context.Context, from string, to string) error {
	st := msgDatabase.Storage

	query := bson.M{
		"$and": []bson.M{
			{"from": from},
			{"to": to},
			{"ischecked": false},
		},
	}

	model := bson.M{
		"ischecked": true,
	}

	return st.Update(ctx, query, model, service.MESSAGE_COLLECTION)
}

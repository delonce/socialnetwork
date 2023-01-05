package friends

import (
	"context"
	"time"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"
)

type FriendManagerService struct {
	friendDatabase FriendQueries
	logger         *logging.Logger
	context        context.Context
}

func NewFriendManager(logger *logging.Logger, config *config.Config) (FriendManager, error) {
	database, err := service.InitNewDatabase(config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &FriendManagerService{
		friendDatabase: NewFriendDB(logger, database),
		logger:         logger,
	}, nil
}

func (manager *FriendManagerService) SendFriendRequest(from string, to string) error {
	newRequest := service.FriendRequest{
		From:       from,
		To:         to,
		IsAccepted: false,
		DateAt:     time.Now(),
	}

	_, err := manager.friendDatabase.AddFriendRequest(manager.context, &newRequest)

	return err
}

func (manager *FriendManagerService) CancelRequest(from string, to string) error {
	fDb := manager.friendDatabase

	return fDb.DeleteRequest(manager.context, from, to)
}

func (manager *FriendManagerService) RejectRequest(from string, to string) error {
	fDb := manager.friendDatabase
	err := fDb.DeleteRequest(manager.context, from, to)

	if err != nil {
		err = fDb.DeleteRequest(manager.context, to, from)

		if err != nil {
			return err
		}
	}

	return nil
}

func (manager *FriendManagerService) AcceptNewFriend(user string, friend string) error {
	requestID, err := manager.GetRequestID(friend, user)

	if err != nil {
		return err
	}

	return manager.friendDatabase.SetAcceptedMark(manager.context, requestID)
}

func (manager *FriendManagerService) DeleteFriend(user string, friend string) error {
	requestID, err := manager.GetRequestID(friend, user)

	if err != nil {
		requestID, err = manager.GetRequestID(user, friend)

		if err != nil {
			return err
		}
	}

	return manager.friendDatabase.SetDeniedMark(manager.context, requestID, friend, user)
}

func (manager *FriendManagerService) GetRequestID(from string, to string) (string, error) {
	request := service.FriendRequest{}

	result, err := manager.friendDatabase.GetRequestByNames(manager.context, from, to)

	if err != nil {
		return "", err
	}

	if err := result.Decode(&request); err != nil {
		manager.logger.Errorf("Error while decoding %s", request)
		return "", err
	}

	return request.ID, nil
}

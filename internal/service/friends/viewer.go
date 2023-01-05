package friends

import (
	"context"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"
)

type FriendViewService struct {
	friendDatabase FriendQueries
	logger         *logging.Logger
	context        context.Context
}

func NewFriendViewer(logger *logging.Logger, config *config.Config) (FriendViewer, error) {
	database, err := service.InitNewDatabase(config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &FriendViewService{
		friendDatabase: NewFriendDB(logger, database),
		logger:         logger,
	}, nil
}

func (viewService *FriendViewService) GetAllProbablyFriends(username string) []service.User {
	userFriends := viewService.GetUserFriends(username)
	userFriends = append(userFriends, username)

	cursor, err := viewService.friendDatabase.GetAllProbFriends(viewService.context, userFriends)

	if err != nil {
		viewService.logger.Panic(err)
	}

	users := []service.User{}
	err = cursor.All(viewService.context, &users)

	if err != nil {
		viewService.logger.Panic(err)
	}

	return users
}

func (viewService *FriendViewService) GetUserFriends(username string) []string {
	cursor, err := viewService.friendDatabase.GetFriends(viewService.context, username)

	if err != nil {
		viewService.logger.Panic(err)
		return nil
	}

	acceptedRequests := []service.FriendRequest{}
	friends := []string{}

	err = cursor.All(viewService.context, &acceptedRequests)

	if err != nil {
		viewService.logger.Panic(err)
		return nil
	}

	for _, accReq := range acceptedRequests {
		if accReq.From != username {
			friends = append(friends, accReq.From)
		} else {
			friends = append(friends, accReq.To)
		}
	}

	return friends
}

func (viewService *FriendViewService) CheckRequest(from string, to string) bool {
	fDb := viewService.friendDatabase

	return fDb.IsExistRequest(viewService.context, from, to)
}

func (viewService *FriendViewService) CheckFriend(from string, to string) bool {
	fDb := viewService.friendDatabase

	return fDb.IsExistFriend(viewService.context, from, to)
}

func (viewService *FriendViewService) GetAllFriendRequests(username string) []service.FriendRequest {
	cursor, err := viewService.friendDatabase.GetAllFriendRequestTo(viewService.context, username)

	if err != nil {
		viewService.logger.Panic(err)
	}

	requests := []service.FriendRequest{}
	err = cursor.All(viewService.context, &requests)

	if err != nil {
		viewService.logger.Panic(err)
	}

	return requests
}

func (viewService *FriendViewService) GetUsersExceptSomeone(exceptUsers []string) []service.User {
	cursor, err := viewService.friendDatabase.GetUsersExcept(viewService.context, exceptUsers)

	if err != nil {
		viewService.logger.Panic(err)
	}

	users := []service.User{}
	err = cursor.All(viewService.context, &users)

	if err != nil {
		viewService.logger.Panic(err)
	}

	return users
}

package friends

import (
	"github.com/delonce/socialnetwork/internal/service"
)

type FriendViewer interface {
	GetAllProbablyFriends(username string) []service.User
	GetAllFriendRequests(username string) []service.FriendRequest
	GetUserFriends(username string) []string

	CheckRequest(from string, to string) bool
	CheckFriend(from string, to string) bool

	// FOR MESSAGE VIEWER
	GetUsersExceptSomeone(exceptUsers []string) []service.User
}

type FriendManager interface {
	GetRequestID(from string, to string) (string, error)

	SendFriendRequest(from string, to string) error
	CancelRequest(from string, to string) error
	RejectRequest(from string, to string) error

	AcceptNewFriend(user string, friend string) error
	DeleteFriend(user string, friend string) error
}

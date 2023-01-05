package messages

import (
	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
)

type MessageManager interface {
	SendMessage(from string, to string, message string) error
	CheckMessage(user string, friend string) error
}

type MessageViewer interface {
	GetOneDialog(from string, to string) []service.ViewMessage
	GetAllDialogs(config *config.Config, username string) []*service.ViewDialog

	CountNewMessages(username string) int64
}

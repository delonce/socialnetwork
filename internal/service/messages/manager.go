package messages

import (
	"context"
	"time"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"
)

type MessageManagerService struct {
	messageDatabase MessageQueries
	logger          *logging.Logger
	context         context.Context
}

func NewMessageManager(logger *logging.Logger, config *config.Config) (MessageManager, error) {
	database, err := service.InitNewDatabase(config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &MessageManagerService{
		messageDatabase: NewMessageDB(logger, database),
		logger:          logger,
	}, nil
}

func (msgManager *MessageManagerService) SendMessage(from string, to string, message string) error {
	newMessage := service.Message{
		From:      from,
		To:        to,
		Text:      message,
		DateAt:    time.Now(),
		IsChecked: false,
	}

	return msgManager.messageDatabase.AddNewMessage(msgManager.context, newMessage)
}

func (msgManager *MessageManagerService) CheckMessage(user string, friend string) error {
	return msgManager.messageDatabase.SetCheckMark(msgManager.context, friend, user)
}

package messages

import (
	"context"
	"sort"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/internal/service/friends"
	"github.com/delonce/socialnetwork/pkg/logging"
)

type MessageViewService struct {
	msgDatabase MessageQueries
	logger      *logging.Logger
	context     context.Context
}

func NewMessageViewer(logger *logging.Logger, config *config.Config) (MessageViewer, error) {
	database, err := service.InitNewDatabase(config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &MessageViewService{
		msgDatabase: NewMessageDB(logger, database),
		logger:      logger,
	}, nil
}

func (msgView *MessageViewService) GetOneDialog(from string, to string) []service.ViewMessage {
	cursor, err := msgView.msgDatabase.GetAllDialogMessages(msgView.context, from, to)

	if err != nil {
		msgView.logger.Panic(err)
	}

	messages := []service.Message{}
	viewMsg := []service.ViewMessage{}

	err = cursor.All(msgView.context, &messages)

	if err != nil {
		msgView.logger.Panic(err)
	}

	for _, msg := range messages {
		viewMsg = append(viewMsg, service.ViewMessage{
			From:       msg.From,
			To:         msg.To,
			Text:       msg.Text,
			FormatDate: msg.DateAt.Format("2006-01-02 15:04"),
			IsChecked:  msg.IsChecked,
		})
	}

	return viewMsg
}

func (msgView *MessageViewService) CountNewMessages(username string) int64 {
	msgAmount := msgView.msgDatabase.CountAllNewMessages(msgView.context, username)

	return msgAmount
}

func (msgView *MessageViewService) GetAllDialogs(config *config.Config, username string) []*service.ViewDialog {
	friendView, err := friends.NewFriendViewer(msgView.logger, config)

	if err != nil {
		msgView.logger.Errorf("Error when creating friend manager, %v", err)
		return nil
	}

	userFriends := friendView.GetUserFriends(username)
	dialogs := []*service.ViewDialog{}

	for _, friend := range userFriends {
		cursor, err := msgView.msgDatabase.GetLastMessage(msgView.context, username, friend)

		if err != nil {
			msgView.logger.Panic(err)
		}

		someDialog := []*service.ViewDialog{}
		err = cursor.All(msgView.context, &someDialog)

		if err != nil {
			msgView.logger.Panic(err)
		}

		if len(someDialog) == 0 {
			dialogs = append(dialogs, &service.ViewDialog{
				From:         friend,
				IsChecked:    true,
				AmountNewMsg: 0,
			})

			continue
		}

		someDialog[0].AmountNewMsg = msgView.msgDatabase.CountNewMessagesBySender(msgView.context, someDialog[0].From, username)

		dialogs = append(dialogs, someDialog[0])
	}

	sort.Slice(dialogs, func(i, j int) (less bool) {
		return dialogs[i].DateAt.Unix() > dialogs[j].DateAt.Unix()
	})

	return dialogs
}

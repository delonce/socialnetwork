package handlers

import (
	"net/http"
	"path"

	"github.com/delonce/socialnetwork/internal/service/friends"
	"github.com/delonce/socialnetwork/internal/service/messages"
	"github.com/julienschmidt/httprouter"
)

func (handler *NetworkHandler) GetMessagePage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := handler.getCurrentUser(w, r)
	friendUsername := params.ByName(USERNAME_URL_TEMPLATE)

	handler.checkConnection(w, r, currentUser.Username, friendUsername)

	msgService, err := messages.NewMessageManager(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Errorf("Error when creating msgService, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	msgService.CheckMessage(currentUser.Username, friendUsername)

	msgViewer, err := messages.NewMessageViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Errorf("Error when creating msgViewer, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	allDialog := msgViewer.GetOneDialog(currentUser.Username, friendUsername)

	templateMap := map[string]interface{}{
		"AllDialog": allDialog,
	}

	SEND_MESSAGE_TEMPLATE.Execute(w, templateMap)
}

func (handler *NetworkHandler) GetMyMessages(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := handler.getCurrentUser(w, r)
	msgViewer, err := messages.NewMessageViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Errorf("Error when creating msgViewer, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	userDialogs := msgViewer.GetAllDialogs(handler.HandlerConfig, currentUser.Username)

	templateMap := map[string]interface{}{
		"CurrentUser": currentUser.Username,
		"Dialogs":     userDialogs,
	}

	ALL_MESSAGES_TEMPLATE.Execute(w, templateMap)
}

func (handler *NetworkHandler) SendNewMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	msgSender := handler.getCurrentUser(w, r)
	msgReciever := params.ByName(USERNAME_URL_TEMPLATE)
	textMessage := r.FormValue("message")

	handler.checkConnection(w, r, msgSender.Username, msgReciever)

	msgService, err := messages.NewMessageManager(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Errorf("Error when creating msgService, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	err = msgService.SendMessage(msgSender.Username, msgReciever, textMessage)

	if err != nil {
		handler.HandlerLogger.Errorf("Error when sending message, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	redirectUrl := path.Join(MESSAGE_URL, msgReciever)

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (handler *NetworkHandler) checkConnection(w http.ResponseWriter, r *http.Request, currentUsername string, friendUsername string) {
	friendView, err := friends.NewFriendViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Errorf("Error when creating friend manager, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	if friendView.CheckFriend(currentUsername, friendUsername) == false {
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}
}

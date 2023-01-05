package handlers

import (
	"net/http"

	"github.com/delonce/socialnetwork/internal/service/friends"
	"github.com/delonce/socialnetwork/internal/service/messages"
	"github.com/delonce/socialnetwork/internal/service/user"

	"github.com/julienschmidt/httprouter"
)

func (handler *NetworkHandler) GetHomePage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := handler.getCurrentUser(w, r)

	friendView, err := friends.NewFriendViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Error when creating friend service, %v", err)
		return
	}

	msgView, err := messages.NewMessageViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Error when creating message service, %v", err)
		return
	}

	otherUsers := friendView.GetAllProbablyFriends(currentUser.Username)
	friendRequests := friendView.GetAllFriendRequests(currentUser.Username)
	msgAmount := msgView.CountNewMessages(currentUser.Username)

	templateMap := map[string]interface{}{
		"CurrentUser":       currentUser,
		"AllUsers":          otherUsers,
		"FriendRequests":    friendRequests,
		"AmountNewMessages": msgAmount,
		"IsCurrentUser":     true,
	}

	HOMEPAGE_TEMPLATE.Execute(w, templateMap)
}

func (handler *NetworkHandler) GetOtherPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := params.ByName(USERNAME_URL_TEMPLATE)
	authService, err := user.NewAuthService(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Can't create auth service, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	friendView, err := friends.NewFriendViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Can't create friend view service, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	otherUser, err := authService.GetUserByName(username)

	if err != nil {
		handler.HandlerLogger.Error("Can't find user, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	currentUser := handler.getCurrentUser(w, r)

	if currentUser.Username == otherUser.Username {
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	templateMap := map[string]interface{}{
		"CurrentUser":          otherUser,
		"IsFriends":            friendView.CheckFriend(currentUser.Username, otherUser.Username),
		"IsReqToPersonExist":   friendView.CheckRequest(currentUser.Username, otherUser.Username),
		"IsReqFromPersonExist": friendView.CheckRequest(otherUser.Username, currentUser.Username),
	}

	HOMEPAGE_TEMPLATE.Execute(w, templateMap)
}

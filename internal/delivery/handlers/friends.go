package handlers

import (
	"net/http"
	"path"

	"github.com/delonce/socialnetwork/internal/service/friends"

	"github.com/julienschmidt/httprouter"
)

func (handler *NetworkHandler) GetFriendRequestsPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := handler.getCurrentUser(w, r)
	friendView, err := friends.NewFriendViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Error when creating friend service, %v", err)
		return
	}

	friendRequests := friendView.GetAllFriendRequests(currentUser.Username)

	templateMap := map[string]interface{}{
		"FriendRequests": friendRequests,
	}

	FRIEND_REQUESTS_TEMPLATE.Execute(w, templateMap)
}

func (handler *NetworkHandler) GetMyFriends(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentUser := handler.getCurrentUser(w, r)
	friendView, err := friends.NewFriendViewer(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Error when creating friend service, %v", err)
		return
	}

	userFriends := friendView.GetUserFriends(currentUser.Username)

	templateMap := map[string]interface{}{
		"Friends": userFriends,
	}

	FRIENDS_TEMPLATE.Execute(w, templateMap)
}

func (handler *NetworkHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	manager, redirectUrl := newManagerURLLink(w, r, *handler, params)

	doFriendRequest(w, r, params, *handler, redirectUrl, manager.SendFriendRequest)
}

func (handler *NetworkHandler) CancelFriendRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	manager, redirectUrl := newManagerURLLink(w, r, *handler, params)

	doFriendRequest(w, r, params, *handler, redirectUrl, manager.CancelRequest)
}

func (handler *NetworkHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	manager, redirectUrl := newManagerURLLink(w, r, *handler, params)

	doFriendRequest(w, r, params, *handler, redirectUrl, manager.RejectRequest)
}

func (handler *NetworkHandler) AcceptFriend(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	manager, redirectUrl := newManagerURLLink(w, r, *handler, params)

	doFriendRequest(w, r, params, *handler, redirectUrl, manager.AcceptNewFriend)
}

func (handler *NetworkHandler) DenyFriend(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	manager, redirectUrl := newManagerURLLink(w, r, *handler, params)

	doFriendRequest(w, r, params, *handler, redirectUrl, manager.DeleteFriend)
}

func newManagerURLLink(w http.ResponseWriter, r *http.Request, handler NetworkHandler, params httprouter.Params) (friends.FriendManager, string) {
	manager, err := friends.NewFriendManager(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Error when creating friend manager, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return nil, ""
	}

	friendUsername := params.ByName(USERNAME_URL_TEMPLATE)
	redirectUrl := path.Join(USERS_URL, friendUsername)

	return manager, redirectUrl
}

func doFriendRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params, handler NetworkHandler,
	nextUrl string, reqFunc func(string, string) error) {

	user := handler.getCurrentUser(w, r)
	friendUsername := params.ByName(USERNAME_URL_TEMPLATE)

	if user.Username == friendUsername {
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	err := reqFunc(user.Username, friendUsername)

	if err != nil {
		handler.HandlerLogger.Error("Error when trying send friend request, %v", err)
		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, nextUrl, http.StatusSeeOther)
}

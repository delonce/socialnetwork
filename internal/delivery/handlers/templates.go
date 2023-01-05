package handlers

import (
	"path"
	"text/template"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

var (
	contextUsername = ""
	ROOT_TEMPLATE   = "web/templates"
	BASE_TEMPLATE   = path.Join(ROOT_TEMPLATE, "base.html")

	INDEX_TEMPLATE           = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "index.html")))
	REGISTER_TEMPLATE        = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "register.html")))
	LOGIN_TEMPLATE           = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "login.html")))
	HOMEPAGE_TEMPLATE        = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "homepage.html"), BASE_TEMPLATE))
	FRIENDS_TEMPLATE         = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "friends.html"), BASE_TEMPLATE))
	FRIEND_REQUESTS_TEMPLATE = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "friend_requests.html"), BASE_TEMPLATE))

	ALL_MESSAGES_TEMPLATE = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "all_messages.html"), BASE_TEMPLATE))
	SEND_MESSAGE_TEMPLATE = template.Must(template.ParseFiles(path.Join(ROOT_TEMPLATE, "send_message.html"), BASE_TEMPLATE))
)

const (
	userIDKey          = "userID"
	accessTokenCookie  = "authAccessToken"
	refreshTokenCookie = "authRefreshToken"
)

type NetworkHandler struct {
	Router        *httprouter.Router
	HandlerLogger *logging.Logger
	HandlerConfig *config.Config
}

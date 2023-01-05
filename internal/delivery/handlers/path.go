package handlers

import "path"

const (
	USERNAME_URL_TEMPLATE = "username"
)

var (
	STATIC_FILE_PATH = "web/static"
	STATIC_FILE_URL  = "/static/*filepath"

	EMPTY_URL             = "/"
	HOME_URL              = "/home"
	REGISTER_URL          = "/register"
	LOGIN_URL             = "/login"
	USERS_URL             = "/users"
	FRIENDS_URL           = "/friends"
	MESSAGE_URL           = "/messages"
	ANY_USERNAME_TEMPLATE = ":" + USERNAME_URL_TEMPLATE

	LOGOUT_URL = path.Join(HOME_URL, "logout")

	OTHER_PAGE_URL = path.Join(USERS_URL, ANY_USERNAME_TEMPLATE)

	FRIEND_REQUESTS_URL = path.Join(FRIENDS_URL, "requests")
	MY_FRIENDS_URL      = path.Join(FRIENDS_URL, "myfriends")

	SEND_REQUEST_URL   = path.Join(USERS_URL, ANY_USERNAME_TEMPLATE, "sendrequest")
	CANCEL_REQUEST_URL = path.Join(USERS_URL, ANY_USERNAME_TEMPLATE, "cancelrequest")
	REJECT_REQUEST_URL = path.Join(USERS_URL, ANY_USERNAME_TEMPLATE, "rejectrequest")

	ACCEPT_FRIEND_URL = path.Join(USERS_URL, ANY_USERNAME_TEMPLATE, "addfriend")
	DENY_FRIEND_URL   = path.Join(USERS_URL, ANY_USERNAME_TEMPLATE, "denyfriend")

	DIALOG_URL = path.Join(MESSAGE_URL, ANY_USERNAME_TEMPLATE)
)

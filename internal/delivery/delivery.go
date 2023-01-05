package delivery

import (
	"net/http"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/delivery/handlers"
	"github.com/delonce/socialnetwork/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

type Delivery interface {
	Register()
	GetRouter() *httprouter.Router
}

type deliveryHandler struct {
	*handlers.NetworkHandler
}

func NewDeliveryHandler(logger *logging.Logger, config *config.Config) *deliveryHandler {
	return &deliveryHandler{
		&handlers.NetworkHandler{
			Router:        httprouter.New(),
			HandlerLogger: logger,
			HandlerConfig: config,
		},
	}
}

func (devHandler *deliveryHandler) Register() {
	devHandler.HandlerLogger.Info("Starting register handlers")

	devHandler.Router.GET(handlers.EMPTY_URL, devHandler.GetIndexPage)

	devHandler.Router.GET(handlers.REGISTER_URL, devHandler.GetRegisterPage)
	devHandler.Router.POST(handlers.REGISTER_URL, devHandler.GoRegisterUser)

	devHandler.Router.GET(handlers.LOGIN_URL, devHandler.GetLoginPage)
	devHandler.Router.POST(handlers.LOGIN_URL, devHandler.SignIn)

	devHandler.Router.GET(handlers.HOME_URL, devHandler.CheckAuth(devHandler.GetHomePage))
	devHandler.Router.GET(handlers.LOGOUT_URL, devHandler.CheckAuth(devHandler.Logout))

	devHandler.Router.GET(handlers.OTHER_PAGE_URL, devHandler.CheckAuth(devHandler.GetOtherPage))
	devHandler.Router.GET(handlers.SEND_REQUEST_URL, devHandler.CheckAuth(devHandler.SendFriendRequest))

	devHandler.Router.GET(handlers.FRIEND_REQUESTS_URL, devHandler.CheckAuth(devHandler.GetFriendRequestsPage))
	devHandler.Router.GET(handlers.MY_FRIENDS_URL, devHandler.CheckAuth(devHandler.GetMyFriends))

	devHandler.Router.GET(handlers.CANCEL_REQUEST_URL, devHandler.CheckAuth(devHandler.CancelFriendRequest))
	devHandler.Router.GET(handlers.REJECT_REQUEST_URL, devHandler.CheckAuth(devHandler.RejectFriendRequest))

	devHandler.Router.GET(handlers.ACCEPT_FRIEND_URL, devHandler.CheckAuth(devHandler.AcceptFriend))
	devHandler.Router.GET(handlers.DENY_FRIEND_URL, devHandler.CheckAuth(devHandler.DenyFriend))

	devHandler.Router.GET(handlers.MESSAGE_URL, devHandler.CheckAuth(devHandler.GetMyMessages))
	devHandler.Router.GET(handlers.DIALOG_URL, devHandler.CheckAuth(devHandler.GetMessagePage))
	devHandler.Router.POST(handlers.DIALOG_URL, devHandler.CheckAuth(devHandler.SendNewMessage))

	devHandler.Router.ServeFiles(handlers.STATIC_FILE_URL, http.Dir(handlers.STATIC_FILE_PATH))

	devHandler.HandlerLogger.Info("Router had registered all handlers")
}

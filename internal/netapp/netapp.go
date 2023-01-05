package netapp

import (
	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/delivery"
	"github.com/delonce/socialnetwork/internal/server"
	"github.com/delonce/socialnetwork/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

type SocialNetwork interface {
	Run()
	startHTTPServer()
}

type NetApp struct {
	NetLogger     *logging.Logger
	AppConfig     *config.Config
	NetworkRouter *httprouter.Router
}

func (networkApp *NetApp) Run() {
	networkApp.startHTTPServer()
}

func InitNewApp(logger *logging.Logger, webConfig *config.Config) *NetApp {
	return &NetApp{
		NetLogger: logger,
		AppConfig: webConfig,
	}
}

func (networkApp *NetApp) getHTTPRouter() *httprouter.Router {
	deliveryHandler := delivery.NewDeliveryHandler(networkApp.NetLogger, networkApp.AppConfig)
	deliveryHandler.Register()

	return deliveryHandler.Router
}

func (networkApp *NetApp) startHTTPServer() {
	networkApp.NetLogger.Info("Settings up router")
	router := networkApp.getHTTPRouter()

	networkApp.NetLogger.Info("Starting http server")
	httpServer := server.GetNewServer(networkApp.AppConfig.Http.Host, networkApp.AppConfig.Http.Port, router)
	networkApp.NetLogger.Infof("HTTP Server had started on %s:%d", networkApp.AppConfig.Http.Host, networkApp.AppConfig.Http.Port)

	if networkApp.AppConfig.Tls.IsTls {
		networkApp.NetLogger.Info("Selected tls mode")
		err := httpServer.ListenAndServeTLS(networkApp.AppConfig.Tls.Cert, networkApp.AppConfig.Tls.Key)

		if err != nil {
			networkApp.NetLogger.Error("Failed to start SSL mode on http server with error %v", err)
		}
	} else {
		networkApp.NetLogger.Info("Selected no tls mode")
		err := httpServer.ListenAndServe()

		if err != nil {
			networkApp.NetLogger.Error("Failed to start noSSL mode on http server with error %v", err)
		}
	}
}

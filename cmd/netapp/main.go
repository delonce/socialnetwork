package main

import (
	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/netapp"
	"github.com/delonce/socialnetwork/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Logger had started")
	webConfig := loadConfig(logger)
	logger.Info("Config had loaded")

	networkApplication := netapp.InitNewApp(logger, webConfig)
	networkApplication.Run()
}

func loadConfig(logger *logging.Logger) *config.Config {
	return config.GetConfig(logger)
}

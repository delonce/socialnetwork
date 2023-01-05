package config

import (
	"sync"

	"github.com/delonce/socialnetwork/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"isDebug"`
	Http    struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
	} `yaml:"http"`

	Tls struct {
		IsTls bool   `yaml:"isTls"`
		Key   string `yaml:"key"`
		Cert  string `yaml:"cert"`
	} `yaml:"tls"`

	Database struct {
		Addr     string `yaml:"address"`
		Port     uint16 `yaml:"port"`
		DBName   string `yaml:"dbname"`
		Password string `yaml:"password"`
		Login    string `yaml:"login"`
	} `yaml:"database"`
}

var instance *Config
var once sync.Once

func GetConfig(logger *logging.Logger) *Config {
	once.Do(func() {
		logger.Info("Getting config from yaml")
		instance = &Config{}
		err := cleanenv.ReadConfig("./configs/config.yaml", instance)

		if err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Error(help)
			logger.Fatal(err)
		}
	})

	logger.Info(instance)
	return instance
}

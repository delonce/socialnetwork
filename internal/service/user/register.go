package user

import (
	"context"
	"fmt"
	"time"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"
)

type RegisterService struct {
	userDatabase UserQueries
	logger       *logging.Logger
	context      context.Context
}

const (
	passwordSalt = "wqfijgopierj889451987*(&9861@!&^5241&@$8513876^1ghhqlkdjhehjgwshjdfuqoo;as"
)

func NewRegisterService(logger *logging.Logger, config *config.Config) (Registration, error) {
	database, err := service.InitNewDatabase(config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &RegisterService{
		userDatabase: NewUserDB(logger, database),
		logger:       logger,
	}, nil
}

func (regServ *RegisterService) RegisterNewUser(username, password, email string) (string, error) {
	newUser := service.User{
		ID:           "",
		Username:     username,
		PasswordHash: getPasswordHash(password),
		Email:        email,
		LastEnt:      time.Now(),
	}

	_, err := regServ.userDatabase.FindUserByUsername(regServ.context, username)
	if err == nil {
		return "", fmt.Errorf("User with username %s already exist!", username)
	}

	_, err = regServ.userDatabase.FindUserByEmail(regServ.context, email)
	if err == nil {
		return "", fmt.Errorf("User with email %s already exist!", email)
	}

	userID, err := regServ.userDatabase.CreateNewUser(regServ.context, &newUser)
	if err != nil {
		regServ.logger.Errorf("Failed to create user %s", newUser)
		return "", err
	}

	return userID, nil
}

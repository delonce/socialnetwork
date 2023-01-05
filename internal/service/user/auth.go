package user

import (
	"context"
	"fmt"

	"github.com/delonce/socialnetwork/internal/config"
	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/pkg/logging"
)

type AuthService struct {
	userDatabase UserQueries
	logger       *logging.Logger
	context      context.Context
}

func NewAuthService(logger *logging.Logger, config *config.Config) (Authorization, error) {
	database, err := service.InitNewDatabase(config)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &AuthService{
		userDatabase: NewUserDB(logger, database),
		logger:       logger,
	}, nil
}

func (auth *AuthService) CreateNewSession(username, password string) (Session, error) {
	user, err := auth.GetUserByCredentials(username, password)

	if err != nil {
		return nil, err
	}

	session, err := NewSession(auth.context, auth.userDatabase, user.ID)

	if err != nil {
		return nil, err
	}

	err = session.DeleteExpiredRefreshTokens()

	if err != nil {
		auth.logger.Errorf("Error when deleting old refresh tokens, user: %s", user.Username)
	}

	_, err = session.RegisterTokenPair()

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (auth *AuthService) CheckSession(accessToken, refreshToken string) (bool, Session) {
	foundedSession, err := FindSession(auth.context, auth.userDatabase, accessToken, refreshToken)

	if err != nil {
		return false, nil
	}

	ok, newSession := foundedSession.SessionIsValid()

	if !ok {
		return false, nil
	}

	if newSession != nil {
		err := foundedSession.ChangeRefreshToken(
			newSession.GetTokenPair().Refresh,
		)

		if err != nil {
			return false, nil
		}

		return true, newSession
	}

	return true, foundedSession
}

func (auth *AuthService) Logout(accessToken, refreshToken string) {
	foundedSession, err := FindSession(auth.context, auth.userDatabase, accessToken, refreshToken)

	if err != nil {
		auth.logger.Errorf("Failed to deleting refresh token %s when logout", refreshToken)
	}

	if err = foundedSession.LogoutRefreshToken(refreshToken); err != nil {
		auth.logger.Errorf("Failed to deleting refresh token %s when logout", refreshToken)
	}
}

func (auth *AuthService) GetUserByCredentials(username, password string) (*service.User, error) {
	foundedUser := service.User{}

	result, err := auth.userDatabase.FindUserByCredentials(auth.context, username, getPasswordHash(password))

	if err != nil {
		return nil, fmt.Errorf("Wrong password or login")
	}

	if err = result.Decode(&foundedUser); err != nil {
		auth.logger.Errorf("Error while decoding %s", username)
		return nil, fmt.Errorf("Server Error")
	}

	return &foundedUser, nil
}

func (auth *AuthService) GetUserByID(userID string) (*service.User, error) {
	foundedUser := service.User{}
	result, err := auth.userDatabase.FindUserByID(auth.context, userID)

	if err != nil {
		auth.logger.Errorf("Can't find user by name %s", userID)
		return nil, err
	}

	if err = result.Decode(&foundedUser); err != nil {
		auth.logger.Errorf("Error while decoding %s", userID)
		return nil, fmt.Errorf("Server Error")
	}

	return &foundedUser, nil
}

func (auth *AuthService) GetUserByName(username string) (*service.User, error) {
	foundedUser := service.User{}
	result, err := auth.userDatabase.FindUserByUsername(auth.context, username)

	if err != nil {
		auth.logger.Errorf("Can't find user by name %s", username)
		return nil, err
	}

	if err = result.Decode(&foundedUser); err != nil {
		auth.logger.Errorf("Error while decoding %s", username)
		return nil, fmt.Errorf("Server Error")
	}

	return &foundedUser, nil
}

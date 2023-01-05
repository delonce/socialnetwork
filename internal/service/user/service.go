package user

import (
	"crypto/sha512"

	"github.com/delonce/socialnetwork/internal/service"
)

type Registration interface {
	RegisterNewUser(username, password, email string) (string, error)
}

type Authorization interface {
	CheckSession(accessToken, refreshToken string) (bool, Session)
	CreateNewSession(username, password string) (Session, error)
	Logout(accessToken, refreshToken string)
	GetUserByCredentials(username, password string) (*service.User, error)
	GetUserByID(userID string) (*service.User, error)
	GetUserByName(username string) (*service.User, error)
}

type Session interface {
	ChangeRefreshToken(newRefresh *service.RefreshToken) error
	DeleteExpiredRefreshTokens() error
	LogoutRefreshToken(refreshTokenUUID string) error
	RegisterTokenPair() (string, error)
	ParseAccessToken() (*tokenClaims, error)
	SessionIsValid() (bool, Session)
	GetTokenPair() *tokenPair
}

func getPasswordHash(password string) string {
	hash := sha512.New()

	hash.Write([]byte(password + passwordSalt))
	bs := hash.Sum(nil)

	return string(bs)
}

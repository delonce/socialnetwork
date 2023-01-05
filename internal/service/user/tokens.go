package user

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/delonce/socialnetwork/internal/service"

	"github.com/dgrijalva/jwt-go"
)

const (
	accessTokenTTL     = 15 * time.Minute
	refreshTokenTTL    = 60 * time.Minute
	refreshTokenLength = 128
	tokenSalt          = "qwd[q;;[]14=-3=359P[qr#$*-!)(#9%*&ijijf[100#)&JJ{Fwwpfefpusdfp83fefgufdst3e"
)

type tokenPair struct {
	Access       string
	Refresh      *service.RefreshToken
	UserDatabase UserQueries
	ctx          context.Context
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID string
}

func NewSession(ctx context.Context, userDatabase UserQueries, tokenUserID string) (Session, error) {
	access, err := GenerateAccessToken(tokenUserID)

	if err != nil {
		return nil, err
	}

	refresh, err := GenerateRefreshToken(tokenUserID)

	if err != nil {
		return nil, err
	}

	return &tokenPair{
		Access:       access,
		Refresh:      refresh,
		UserDatabase: userDatabase,
		ctx:          ctx,
	}, nil
}

func FindSession(ctx context.Context, db UserQueries, accessToken string, refreshToken string) (Session, error) {
	foundedRefresh := service.RefreshToken{}
	result, err := db.FindRefreshTokenByUUID(ctx, refreshToken)

	if err != nil {
		return nil, fmt.Errorf("Wrong password or login")
	}

	if err = result.Decode(&foundedRefresh); err != nil {
		return nil, fmt.Errorf("Error when decoding refresh token")
	}

	return &tokenPair{
		Access:       accessToken,
		Refresh:      &foundedRefresh,
		UserDatabase: db,
		ctx:          ctx,
	}, nil
}

func (pair *tokenPair) GetTokenPair() *tokenPair {
	return pair
}

func (pair *tokenPair) SessionIsValid() (bool, Session) {
	_, err := pair.ParseAccessToken()

	if err != nil {
		if time.Now().Unix() > pair.Refresh.ExpiresAt.Unix() {
			err := pair.UserDatabase.DeleteRefreshToken(pair.ctx, pair.Refresh.UUID)

			if err != nil {
				panic(err)
			}

			return false, nil

		} else {

			newRefreshToken, err := GenerateRefreshToken(pair.Refresh.UserID)
			newRefreshToken.ExpiresAt = pair.Refresh.ExpiresAt

			if err != nil {
				return false, nil
			}

			newAccessToken, err := GenerateAccessToken(newRefreshToken.UserID)

			if err != nil {
				return false, nil
			}

			return true, &tokenPair{
				Access:       newAccessToken,
				Refresh:      newRefreshToken,
				UserDatabase: pair.UserDatabase,
				ctx:          pair.ctx,
			}
		}
	}

	return true, nil
}

func (pair *tokenPair) ChangeRefreshToken(newRefresh *service.RefreshToken) error {
	err := pair.UserDatabase.DeleteRefreshToken(pair.ctx, pair.Refresh.UUID)

	if err != nil {
		return err
	}

	_, err = pair.UserDatabase.AddRefreshToken(pair.ctx, newRefresh)

	if err != nil {
		return err
	}

	return nil
}

func (pair *tokenPair) RegisterTokenPair() (string, error) {
	sessionID, err := pair.UserDatabase.AddRefreshToken(pair.ctx, pair.Refresh)

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (pair *tokenPair) ParseAccessToken() (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(pair.Access, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Token has invalid signing method")
		}

		return []byte(tokenSalt), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return nil, fmt.Errorf("Invalid type for decoding token")
	}

	return claims, nil
}

func (pair *tokenPair) DeleteExpiredRefreshTokens() error {
	cursor, _ := pair.UserDatabase.FindExpiredRefreshTokens(
		pair.ctx,
		pair.Refresh.UserID,
	)

	var tokens []service.RefreshToken

	if err := cursor.All(context.TODO(), &tokens); err != nil {
		return err
	}

	for _, token := range tokens {
		err := pair.UserDatabase.DeleteRefreshToken(
			pair.ctx,
			token.UUID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (pair *tokenPair) LogoutRefreshToken(refreshTokenUUID string) error {
	if err := pair.UserDatabase.DeleteRefreshToken(pair.ctx, refreshTokenUUID); err != nil {
		return err
	}

	return nil
}

func GenerateAccessToken(tokenUserID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},

		tokenUserID,
	})

	return token.SignedString([]byte(tokenSalt))
}

func GenerateRefreshToken(tokenUserID string) (*service.RefreshToken, error) {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	var builder strings.Builder

	for i := 0; i < refreshTokenLength; i++ {
		builder.WriteRune(chars[rand.Intn(len(chars))])
	}

	tokenUUID := builder.String()

	return &service.RefreshToken{
		ID:        "",
		UUID:      tokenUUID,
		UserID:    tokenUserID,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}, nil
}

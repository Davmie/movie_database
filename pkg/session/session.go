package session

import (
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/pkg/errors"
)

var tokenKey = []byte("fvoNImvpdms023sv0s9vs")

type UserClaims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

type Claims struct {
	User UserClaims `json:"user"`
	jwt.RegisteredClaims
}

type JWTSessionsManager struct{}

func (jsm JWTSessionsManager) GetUser(inToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(inToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenKey, nil
	})

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return -1, "", errors.Wrapf(err, "can`t parse or validate session token \"%s\"", inToken)
	}

	return claims.User.ID, claims.User.Role, nil
}

func (jsm JWTSessionsManager) CreateSession(id int, role string) (string, error) {
	claims := Claims{
		UserClaims{
			ID:   id,
			Role: role,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(tokenKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert token to string")
	}

	return tokenString, nil
}

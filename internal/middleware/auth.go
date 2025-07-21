package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidJWT = errors.New("invalid JWT token , authorize")
)

func GenerateJWTToken(login string, JWTSecret string) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 5).Unix(), // время действия токена
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// подпись токена
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", fmt.Errorf("[GenerateJWTToken| sign token]: %w", err)
	}
	return tokenString, nil
}

func checkTokenIsValid(tokenString, JWTsecret string) error {
	if tokenString != "" {
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(JWTsecret), nil
		})

		if err != nil || !token.Valid {
			return ErrInvalidJWT
		}

	}
	return nil
}

func getLoginFromValidToken(tokenString, JWTsecret string) interface{} {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTsecret), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims["login"].(string)
}

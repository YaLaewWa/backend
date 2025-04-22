package util

import (
	"fmt"
	"socket/pkg/apperror"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
	JWTSecretKey string `env:"SECRET,required"`
}

type JWTUtils struct {
	Config *JWTConfig
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func NewJWTUtils(config *JWTConfig) *JWTUtils {
	return &JWTUtils{Config: config}
}

func (j *JWTUtils) GenerateJWT(username string) (string, error) {
	claims := &JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(j.Config.JWTSecretKey))
	if err != nil {
		return "", apperror.InternalServerError(err, "cannot generate token")
	}
	return tokenString, nil
}

func (j *JWTUtils) DecodeJWT(inputToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Config.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return new(JWTClaims), apperror.UnauthorizedError(err, "parse token failed")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return new(JWTClaims), apperror.UnauthorizedError(err, "invalid token")
	}

	return claims, nil
}

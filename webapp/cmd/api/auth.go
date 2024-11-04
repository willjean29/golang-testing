package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtTokenExpiry = time.Minute * 15
const jwtRefreshTokenExpiry = time.Hour * 24

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Clams struct {
	Username string `json:"name"`
	jwt.RegisteredClaims
}

func (app *application) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Clams, error) {
	w.Header().Add("Vary", "Authorization")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("no authorization header")
	}

	headersParts := strings.Split(authHeader, " ")
	if len(headersParts) != 2 {
		return "", nil, errors.New("invalid authorization header")
	}

	if headersParts[0] != "Bearer" {
		return "", nil, errors.New("invalid authorization header")
	}

	token := headersParts[1]

	claims := &Clams{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(app.JWTSecret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("token is expired")
		}
		return "", nil, errors.New("invalid token")
	}

	if claims.Issuer != app.Domain {
		return "", nil, errors.New("incorrect issuer")
	}

	return token, claims, nil
}

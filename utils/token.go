package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(id int) (string, error) {
	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("SECRET_KEY"))
	if err != nil {
		return "", err
	}
	return token, err
}

func ValidateToken(signedToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("SECRET_KEY"), nil
		},
	)

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return 0, errors.New("token is expired")
	}

	return claims.ID, nil
}

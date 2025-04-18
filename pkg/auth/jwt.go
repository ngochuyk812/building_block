package auth_context

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimModel struct {
	IdSite     string   `json:"id_site"`
	IdAuthUser string   `json:"id_auth_user"`
	Roles      []string `json:"roles"`
	UserName   string   `json:"user_name"`
	Email      string   `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(claims *ClaimModel, serect string, expireIn time.Duration) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(serect)
}

func VerifyJWT(tokenStr string, serect string) (*ClaimModel, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimModel{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return serect, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ClaimModel); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

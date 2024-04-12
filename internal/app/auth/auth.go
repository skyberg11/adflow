package auth

import (
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

var JwtKey = []byte("DEFUALT_SECRET")

func GenerateJWT(id int64) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtKey)
}

func ValidateToken(reqToken string, id int64) (int, string) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, "unauthorized"
		}
		return http.StatusBadRequest, "bad request"
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, "unauthorized"
	}

	if claims.ID != id {
		return http.StatusForbidden, "unmatch"
	}

	return 200, "ok"
}

package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)





func makeJWT(userId uuid.UUID, tokenSecret string, expiresAt time.Duration) (string, error) { 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer: "smartDesigner-access",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		Subject: userId.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err 
	}

	return signedToken, nil 
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil 
	})
	if err != nil || !token.Valid{
		return uuid.Nil, err 
	}
	stringfiedUserID, err := token.Claims.GetSubject()
	userId, err := uuid.Parse(stringfiedUserID)
	if err != nil {
		log.Fatalf("cann't pares the stringified user id to an uuid \n%s", err)
	}
	return userId, nil 
}




package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var API_SECRET = "wP3-sN6&gG4-lV8>gJ9)"
var token_lifespan = "100"

func GenerateToken(user_id uint) (string, error) {

	token_lifespan, err := strconv.Atoi(token_lifespan)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(API_SECRET))
}

func TokenValid(tokenString string) error {
	
	fmt.Println(tokenString)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	})
	if err != nil {
		return err
	}

	return nil
}

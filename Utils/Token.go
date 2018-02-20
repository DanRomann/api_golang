package Utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
	"fmt"
	"errors"
)

var secret = "haha-secret"


func CreateToken(userId int, day int) (*string ,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(day)).Unix(),
		"userId": userId,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil{
		log.Println("Utils.Token.CreateToken ", err)
		return nil, errors.New("can't create token")
	}
	return &tokenString, nil

}

func ParseToken(tokenString *string) (int, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["userId"].(float64); ok{
			return int(userId), nil
		}
		return 0, errors.New("rotten token")

	} else {
		return 0, nil
	}
}
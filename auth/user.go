package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type User struct {
	Username string
	Email string
	Password string
}


func prepareToken(user *User) string {
	tokenContent := jwt.MapClaims{
		"username": user.Username,
		"expiry": time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	HandleErr(err)

	return token
}
package auth

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type validation struct {
	Valid string
	Value string
}

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)
	return string(hashed)
}

// Create validation
func Validation(values []validation) bool {
	username := regexp.MustCompile(`^([A-Za-z0-9]{5,})+$`)
	//email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$`)

	for _, value := range values {
		switch value.Valid {
		case "username":
			if !username.MatchString(value.Value) {
				return false
			}
		case "email":
			if len(value.Value) < 8 {
				return false
			}
		case "password":
			if len(value.Value) < 5 {
				return false
			}
		}
	}
	return true
}

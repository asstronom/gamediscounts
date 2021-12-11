package auth

import (
	"encoding/json"
	"fmt"
	_ "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var DB = map[string]string{}

type Credentials struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp (w http.ResponseWriter, r *http.Request){
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader ( http.StatusBadRequest)
		return
	}
	if !Validation([]validation{{"username",creds.Username},{"email",creds.Email}, {"password",creds.Password}}){
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":"some of filled values is not valid"})
		return
	}
	hashedPassword ,_:= bcrypt.GenerateFromPassword([]byte(creds.Password),8) // 8 is a cost of hashing (some specific crypto shit)

	//<----------------USER DB INSERTION CODE GOES HERE -------------------->
	DB[creds.Username] = string(hashedPassword) // testing map instead of DB

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(DB)
}

func SignIn ( w http.ResponseWriter , r *http.Request){
	creds := &Credentials{}
	err:= json.NewDecoder(r.Body).Decode(creds)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !Validation([]validation{{"username",creds.Username}, {"password",creds.Password}}){
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":"some of filled values is not valid"})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(DB[creds.Username]), []byte(creds.Password)); err != nil { // replace map with DB
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":"wrong username or password"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":"Welcome",
	})
}

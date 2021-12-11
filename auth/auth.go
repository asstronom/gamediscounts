package auth

import (
	"encoding/json"
	"fmt"
	userdb "github.com/gamediscounts/db/couchdb"
	_ "github.com/golang-jwt/jwt"
	"github.com/leesper/couchdb-golang"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

const (
	userDatabaseURL  = "http://couchdb:couchdb@localhost:5984"
	userDatabaseName = "gamediscounts"
)

//var DB = map[string]string{}

type Credentials struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !Validation([]validation{{"username", creds.Username}, {"email", creds.Email}, {"password", creds.Password}}) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "some of filled values is not valid"})
		return
	}
	emailParts := strings.Split(creds.Email, "@")
	for _, item := range emailParts {
		fmt.Println(item)
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 8) // 8 is a cost of hashing (some specific crypto shit)
	//newUser := userdb.User{Credentials: userdb.Credentials{ creds.Username, string(hashedPassword)}, Document: couchdb.Document{}}
	newUser := userdb.User{userdb.Credentials{creds.Username, string(hashedPassword)}, emailParts[0], emailParts[1], false, false, false, couchdb.Document{}}
	fmt.Println(newUser)
	//<----------------USER DB INSERTION CODE GOES HERE -------------------->
	//DB[creds.Username] = string(hashedPassword) // testing map instead of DB
	//db, err := userdb.OpenDB(userDatabaseURL, userDatabaseName)
	db, err := userdb.OpenDB("http://couchdb:couchdb@localhost:5984", "gamediscounts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	_, err = db.AddUser(newUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	//fmt.Println(id)

	//fmt.Println(DB)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !Validation([]validation{{"username", creds.Username}, {"password", creds.Password}}) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "some of filled values is not valid"})
		return
	}
	db, err := userdb.OpenDB(userDatabaseURL, userDatabaseName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	user, err := db.GetUserByName(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		log.Println(err)
		return
	}
	hashedPassword := user.HashedPassword

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil { // replace map with DB
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "wrong username or password"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome",
	})

	//TODO
	//add JW Token
}

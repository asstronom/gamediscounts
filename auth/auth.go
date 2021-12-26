package auth

import (
	"encoding/json"
	"fmt"
	userdb "github.com/gamediscounts/db/couchdb"
	wishlist "github.com/gamediscounts/db/neo4j"
	"github.com/golang-jwt/jwt"
	"github.com/leesper/couchdb-golang"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	userDatabaseURL  = "http://couchdb:couchdb@localhost:5984"
	userDatabaseName = "gamediscounts"
	jwtKey           = "sheesh"
)

//var DB = map[string]string{}

type Credentials struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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
	userDB, err := userdb.OpenDB("http://couchdb:couchdb@localhost:5984", "gamediscounts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	_, er := userDB.AddUser(newUser)
	if er != nil {
		log.Fatalln("error opening wishlist: ", er)
		return
	}
	wishlistDB, er := wishlist.OpenDB("neo4j://localhost:7687", "neo4j", "GuesgP4LPLS")
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = wishlistDB.AddUser(newUser.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Method", "true")
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
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte("sheesh")) // should use env variable for secret
	if err != nil {
		log.Println(err)
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	fmt.Println(tokenString)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome",
	})

}
func GetTokenUsername(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		log.Println(err)
	}
	claims := &Claims{}
	tknStr := c.Value

	_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	return claims.Username, err
}
func FetchUserName(w http.ResponseWriter, r *http.Request) {
	username, err := GetTokenUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"username": username})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Method", "true")
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				json.NewEncoder(w).Encode(nil) // iotii is fucking trash frond end dev
				log.Println(err)
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		claims := &Claims{}
		tknStr := c.Value
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil {
			log.Println(err)
			if err == jwt.ErrSignatureInvalid {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			log.Println("Token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			fmt.Println("userName:", claims.Username)
			endpoint(w, r)
		}
	}
}
func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})
}

package userdb

import (
	"fmt"

	"github.com/leesper/couchdb-golang"
)

type UserDB struct {
	server *couchdb.Server
	db     *couchdb.Database
}

type Credentials struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashedPassword"`
}

type User struct {
	Credentials
	EmailName               string `json:"emailname"`
	EmailDomain             string `json:"emaildomain"`
	Verified                bool   `json:"verified"`
	SubscribedEmailWishlist bool   `json:"subscribedEmailWishlist"`
	SubscribedPushWishlist  bool   `json:"subscribedPushWishlist"`
	couchdb.Document
}

func OpenDB(urlstr string, dbname string) (UserDB, error) {
	var resDB UserDB
	server, err := couchdb.NewServer(urlstr)
	if err != nil {
		fmt.Println("Error connecting to server")
		return UserDB{}, err
	}

	db, err := server.Get(dbname)

	if err != nil {
		fmt.Println("Error connecting to DB")
		return UserDB{}, err
	}

	resDB.server = server
	resDB.db = db

	return resDB, nil
}

func (db *UserDB) CheckUsernameExists(username string) (bool, error) {
	selector := fmt.Sprintf("username == %s", username)
	doc, err := db.db.Query(nil, selector, nil, nil, nil, []string{"registration", "username"})
	if err != nil {
		return true, err
	}
	if len(doc) != 0 {
		return true, nil
	}

	return false, nil
}

func (db *UserDB) CheckEmailExists(emailName string, emailDomain string) (bool, error) {
	selector := fmt.Sprintf("emailname == %s && emaildomain == %s", emailName, emailDomain)
	doc, err := db.db.Query(nil, selector, nil, nil, nil, []string{"registration", "username"})
	if err != nil {
		return true, err
	}
	if len(doc) != 0 {
		return true, nil
	}

	return false, nil
}

func (db *UserDB) CheckIfCredentialsAreValid(user User) (bool, error) {
	b, err := db.CheckUsernameExists(user.Username)
	if err != nil {
		return false, err
	}
	if b {
		return false, fmt.Errorf("username already exists")
	}

	b, err = db.CheckEmailExists(user.EmailName, user.EmailDomain)
	if err != nil {
		return false, err
	}
	if b {
		return false, fmt.Errorf("email already exists")
	}

	return true, nil
}

func (db *UserDB) AddUser(user User) (string, error) {

	b, err := db.CheckIfCredentialsAreValid(user)

	if err != nil {
		return "", err
	}

	if !b {
		return "", err
	}

	var docid string
	for {
		docid = couchdb.GenerateUUID()
		if db.db.Contains(docid) != nil {
			break
		}
	}

	err = user.SetID(docid)
	if err != nil {
		return "", err
	}

	doc, err := couchdb.ToJSONCompatibleMap(user)
	if err != nil {
		fmt.Println("Error converting to JSONCompatibleMap", err)
		return "", err
	}

	err = db.db.Set(docid, doc)
	if err != nil {
		fmt.Println("Error setting document")
		return "", err
	}

	return docid, nil
}

func (db *UserDB) GetUserByName(username string) (User, error) {
	selector := fmt.Sprintf("username == %s", username)
	doc, err := db.db.Query(nil, selector, nil, nil, nil, []string{"_design/registration", "username"})
	if err != nil {
		return User{}, err
	}
	if len(doc) != 1 {
		return User{}, fmt.Errorf("error while searching")
	}
	user := User{}
	err = couchdb.FromJSONCompatibleMap(&user, doc[0])
	if err != nil {
		fmt.Println("error converting to struct User")
		return User{}, err
	}

	return user, nil
}

func (db *UserDB) GetUserByEmail(emailName string, emailDomain string) (User, error) {
	selector := fmt.Sprintf("emailname == %s && emaildomain == %s", emailName, emailDomain)
	doc, err := db.db.Query(nil, selector, nil, nil, nil, []string{"registration", "email"})
	if err != nil {
		fmt.Println("Here")
		return User{}, err
	}
	if len(doc) != 1 {
		return User{}, fmt.Errorf("error while searching")
	}
	user := User{}
	err = couchdb.FromJSONCompatibleMap(&user, doc[0])
	if err != nil {
		fmt.Println("error converting to struct User")
		return User{}, err
	}

	return user, nil
}

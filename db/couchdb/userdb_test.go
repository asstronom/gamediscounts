package userdb

import (
	"fmt"
	"testing"

	"github.com/leesper/couchdb-golang"
)

func TestAddUser(t *testing.T) {
	db, err := OpenDB("http://couchdb:couchdb@localhost:5984", "test")
	if err != nil {
		t.Errorf("Error opening database")
	}

	testTable := []struct {
		user        User
		expectedErr error
	}{
		{
			User{Credentials{"asstronom", "secret"}, "danya.live", "gmail.com", false, false, false, couchdb.Document{}},
			nil,
		},
		{
			User{Credentials{"asstronom", "secret"}, "danya.live", "gmail.com", false, false, false, couchdb.Document{}},
			fmt.Errorf("username already exists"),
		},
		{
			User{Credentials{"asstronomer", "secret"}, "danya.live", "gmail.com", false, false, false, couchdb.Document{}},
			fmt.Errorf("email already exists"),
		},
	}
	for i, v := range testTable {
		_, err := db.AddUser(v.user)
		if err != v.expectedErr {
			if err.Error() != v.expectedErr.Error() {
				t.Errorf("Incorrect error %d, %s != %s", i, err.Error(), v.expectedErr.Error())
			}
		}
	}

	for _, v := range testTable {
		db.RemoveUserByName(v.user.Username)
	}
	len, err := db.db.Len()
	if err != nil {
		t.Errorf("Error while cleaning up 1")
	}
	if len != 0 {
		t.Errorf("Error while cleaning up 2")
	}
}

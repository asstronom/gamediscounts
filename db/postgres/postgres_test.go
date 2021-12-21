package postgres

import (
	"fmt"
	"testing"

	"github.com/gamediscounts/model/steamapi"
)

const (
	host     = "localhost"
	port     = 5432
	username = "user"
	password = "mypassword"
	dbname   = "gamediscounts"
)

func TestInsertDLC(t *testing.T) {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := Open(postgresInfo) // dummy DB for test
	if err != nil {
		t.Errorf(err.Error())
	}
	app, err := steamapi.GetAppInfo(378648, "ua")
	if err != nil {
		t.Errorf("error while extracting appinfo")
	}
	err = db.insertDLC(app, -1)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestInsertPackage(t *testing.T) {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := Open(postgresInfo) // dummy DB for test
	if err != nil {
		t.Errorf(err.Error())
	}
	pack, err := steamapi.GetPackageInfo(124923, "ua")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = db.insertPackage(pack)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestMatchGenre(t *testing.T) {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := Open(postgresInfo) // dummy DB for test
	if err != nil {
		t.Errorf(err.Error())
	}
	res, err := db.matchGenre(`Action`)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf(`Action: %d`, res)
	res, err = db.matchGenre(`Action`)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf(`Action: %d`, res)
	res, err = db.matchGenre(`Adventure`)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf(`Action: %d`, res)
}

func TestInsertGame(t *testing.T) {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := Open(postgresInfo) // dummy DB for test
	if err != nil {
		t.Errorf(err.Error())
	}
	appinfo, err := steamapi.GetAppInfo(292030, "ua")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = db.insertGame(appinfo)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestInitDB(t *testing.T) {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := Open(postgresInfo) // dummy DB for test
	if err != nil {
		t.Errorf(err.Error())
	}
	err = db.InitDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
}

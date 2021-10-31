package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/tidwall/gjson"
)

const (
	host     = "postgres"
	port     = 5432
	username = "user"
	password = "mypassword"
	dbname   = "gamediscounts"
)

// Get this package if it's missing.
// go get -u github.com/lib/p/ go get -u github.com/lib/pq

func Open(credentials string) *sql.DB {
	db, err := sql.Open("postgres", credentials)

	if err != nil {
		log.Fatalln(err)
	}

	start := time.Now()
	for db.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			fmt.Println("failed to connect after 10 secs.")
			break
		}
	}
	fmt.Println("connected:", db.Ping() == nil)

	return db
}

func main() {
	fmt.Println("connecting")
	// these details match the docker-compose.yml file.
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db := Open(postgresInfo)
	defer db.Close()

	// _, err := db.Exec(`DROP TABLE IF EXISTS game;`)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	res := gjson.Get(string(body), "applist.apps").Array()
	output := fmt.Sprintf("Number of games: %d", len(res))
	fmt.Println(output)
	fmt.Println(res[0].Get("appid"), res[0].Get("name"))

	sqlQuery := fmt.Sprintf(`INSERT INTO game(name) VALUES ('%s')`, res[0].Get("name").Value())
	_, err = db.Exec(sqlQuery)

	if err != nil {
		log.Fatalln(err)
	}

	// _, err = db.Exec(`CREATE TABLE COMPANY (ID INT PRIMARY KEY NOT NULL, NAME text);`)
	// if err != nil {
	// 	panic(err)
	// }
	//fmt.Println("table company is created")
}

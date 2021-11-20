package main

import (
	"database/sql"
	"fmt"
	"github.com/gamediscounts/db/postgres"

	//"io/ioutil"
	"log"
	//"net/http"

	_ "github.com/lib/pq"
	//"github.com/tidwall/gjson"
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

func run() {
	fmt.Println("connecting")
	// these details match the docker-compose.yml file.
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	// resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// res := steamapi.GetAppList()
	// output := fmt.Sprintf("Number of games: %d", len(res))
	// fmt.Println(output)
	// fmt.Println(res[0].Get("appid"), res[0].Get("name"))

	//err = db.InitGames()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(steamapi.GetAppPrice(271590, "ua"))
	// appids := []int{620, 1203220, 271590, 216938}
	// fmt.Println(steamapi.GetAppsPrice(&appids, "ua"))
	//db.TestQueryRow()
	//err = db.TestQuery()
	if err != nil {
		log.Fatal("WHY", err)
	}
	fmt.Println("Initializing")
	err = db.InitGamePrice()
	if err != nil {
		log.Fatal(err)
	}
	//sqlQuery := fmt.Sprintf(`INSERT INTO game(name) VALUES ('%s')`, res[0].Get("name").Value())
	//_, err = db.Exec(sqlQuery)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// _, err = db.Exec(`CREATE TABLE COMPANY (ID INT PRIMARY KEY NOT NULL, NAME text);`)
	// if err != nil {
	// 	panic(err)
	// }
	//fmt.Println("table company is created")
}

func main() {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, err := sql.Open("postgres", postgresInfo)

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Query("SELECT storegameid, gameid FROM gameprice")

	if err != nil {
		log.Fatalln(err)
	}

	if !rows.Next() {
		log.Fatalln("FUCK THIS SHIT IM OUT")
	}
	var res1 int
	var res2 int
	err = rows.Scan(&res1, &res2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res1, res2)

	if !rows.Next() {
		log.Fatalln("FUCK THIS SHIT IM OUT")
	}

	err = rows.Scan(&res1, &res2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res1, res2)

	solvedata, err := postgres.OpenSolve(postgresInfo)

	if err != nil {
		log.Fatalln(err)
	}

	solvedata.SolveQuery()

	//run()

	if err != nil {
		log.Fatalln(err)
	}
}

package main

import (

	"context"

	"fmt"
	"github.com/gamediscounts/db/postgres"
	"github.com/gamediscounts/server"
	"net/http"
	"time"

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

func run() error {
	fmt.Println("connecting")
	// these details match the docker-compose.yml file.
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)

	if err != nil {
		return err
	}

	defer db.Close()

	// err = db.InitTables()

	// if err != nil {
	// 	return err
	// }

	// err = db.InitStores()
	// if err != nil {
	// 	return err
	// }

	// err = db.InitGames()
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
	return nil
}

func main() {

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
	var solvedata *postgres.SolveDB // dummy DB for test

	ctx:= context.Background()
	s:= server.Init(ctx,solvedata)
	addr:=":8080"

	httpServer:= &http.Server{
		Addr: addr,
		Handler: s,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("staring web server on %s",addr)
	if err := httpServer.ListenAndServe(); err!=nil{
		log.Fatalln( err)
	}
}


package main

import (
	"context"

	"fmt"
	"net/http"
	"time"

	"github.com/gamediscounts/db/postgres"
	"github.com/gamediscounts/server"

	//"io/ioutil"
	"log"
	//"net/http"

	_ "github.com/lib/pq"
	//"github.com/tidwall/gjson"
)

const (
	host     = "localhost"
	port     = 5432
	username = "user"
	password = "mypassword"
	dbname   = "gamediscounts"
)

// Get this package if it's missing.
// go get -u github.com/lib/p/ go get -u github.com/lib/pq

func initdb() error {
	fmt.Println("connecting")
	// these details match the docker-compose.yml file.
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.InitTables()

	if err != nil {
		return err
	}

	fmt.Println("Inited tables")

	err = db.InitStores()
	if err != nil {
		return err
	}

	err = db.InitGames()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initializing")
	err = db.InitGamePrice()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

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

	err = db.RefreshFeatured()
	if err != nil {
		log.Fatalln(err)
	}

	res, err := db.BestOffers(0, 8, postgres.UA)
	if err != nil {
		log.Fatalln()
	}

	res1, err := db.GetGame(126074, postgres.UA)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
	fmt.Println(len(res))
	fmt.Println(res1)

	return nil
}

func main() {

	//errInit := initdb()
	//if errInit != nil {
	//	log.Fatalln(errInit)
	//}

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("connecting")
	// these details match the docker-compose.yml file.
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo) // dummy DB for test

	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	s := server.Init(ctx, db)
	addr := ":8080"

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("staring web server on %s", addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

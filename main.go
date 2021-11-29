package main

import (
	"fmt"

	//"io/ioutil"
	"log"
	//"net/http"

	"github.com/gamediscounts/postgres"
	"github.com/gamediscounts/steamapi"
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

func TestGetAppPrice() error {
	fmt.Println("GetAppPrice")
	fmt.Println(steamapi.GetAppPrice(620, "ua"))
	fmt.Println(steamapi.GetAppPrice(570, "ua"))
	fmt.Println(steamapi.GetAppPrice(271590, "ua"))
	fmt.Println(steamapi.GetAppPrice(400, "ua"))
	fmt.Println(steamapi.GetAppPrice(216938, "ua"))
	//fmt.Println(steamapi.GetAppPrice(620, "ua"))
	//fmt.Println(steamapi.GetAppPrice(620, "ua"))
	return nil
}

func TestGetAppsPrice() error {
	fmt.Println("GetAppsPrice")
	fmt.Println(steamapi.GetAppsPrice(&[]int{620, 570, 271590, 400, 216938}, "ua"))
	return nil
}

func TestGetFeaturedCategories() error {
	res1, res2, err := steamapi.GetFeaturedCategories("ua")
	if err != nil {
		return err
	}
	fmt.Println("GetFeaturedCategories\n", res1, res2)
	return nil
}

func TestRefreshFeaturedCategories() error {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)
	if err != nil {
		return err
	}
	err = db.RefreshFeatured()
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func TestGetBestPrice() error {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)
	if err != nil {
		return err
	}

	fmt.Println(db.BestOffers("ua"))
	return nil
}

func TestGetAppPriceDB() error {
	fmt.Println("TestGetAppPriceDB")
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)
	if err != nil {
		return err
	}
	resid, resstores, err := db.BestOffers("ua")
	if err != nil {
		return err
	}
	for i := 0; i < len(resid); i++ {
		fmt.Println(db.GetAppPrice(resid[i], resstores[i], "ua"))
	}
	return nil
}

func TestTGetGameName() error {
	fmt.Println("TestGetGameName")
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := postgres.Open(postgresInfo)
	if err != nil {
		return err
	}
	resid, _, err := db.BestOffers("ua")
	if err != nil {
		return err
	}
	for i := 0; i < len(resid); i++ {
		fmt.Println(db.GetGameName(resid[i]))
	}
	return nil
}

func RunTests() error {
	TestGetAppPrice()
	TestGetAppsPrice()
	TestGetFeaturedCategories()
	TestRefreshFeaturedCategories()
	TestGetBestPrice()
	TestGetAppPriceDB()
	TestTGetGameName()
	return nil
}
func main() {
	// err := run()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	RunTests()
}

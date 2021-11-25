package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gamediscounts/steamapi"
)

var (
	steamid int = 1
)

type GameDB struct {
	*sql.DB
}

type GamePrice struct {
	storegameid string
	price       int
	discount    int
	isFree      bool
}

func Open(credentials string) (*GameDB, error) {
	db, err := sql.Open("postgres", credentials)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	for db.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			return nil, err
		}
	}
	fmt.Println("connected:", db.Ping() == nil)
	database := GameDB{db}
	return &database, nil
}

func (DB *GameDB) CloseDB() {
	DB.Close()
}

func (DB *GameDB) InitTables() error {
	// _, err := DB.Exec("DROP DATABASE IF EXISTS gamediscounts")
	// if err != nil {
	// 	return err
	// }
	_, err := DB.Exec("DROP TABLE IF EXISTS game CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS gameprice CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS store CASCADE")
	if err != nil {
		return err
	}
	// _, err = DB.Exec("CREATE DATABASE gamediscounts")
	// if err != nil {
	// 	return err
	// }
	_, err = DB.Exec(`CREATE TABLE game (
		id SERIAL PRIMARY KEY,
		name varchar(255) UNIQUE)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE TABLE store (
		id SERIAL PRIMARY KEY,
		name varchar(255) UNIQUE)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE TABLE gameprice (
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE,
		CONSTRAINT gamePriceId PRIMARY KEY (gameid, storeid),
		storegameid VARCHAR(255) UNIQUE,
		price NUMERIC,
		discount INT DEFAULT 0, 
		free BOOLEAN DEFAULT FALSE)`)
	if err != nil {
		return err
	}
	return nil
}

func (DB *GameDB) InitStores() error {
	//query := fmt.Sprintf(`INSERT INTO store (id, name) VALUES (%d, 'steam')`, steamid)
	_, err := DB.Exec(`INSERT INTO store (id, name) VALUES ($1, 'steam')`, steamid)
	if err != nil {
		fmt.Println("Error InitStores")
		return err
	}

	return nil
}

func (DB *GameDB) InitGames() error {
	res := steamapi.GetAppList()
	for i := 0; i < len(res); i++ {
		var curSteamId int = int(res[i].Get("appid").Value().(float64))
		var curName string = res[i].Get("name").Value().(string)
		//fmt.Println("Current game: ", curName, curSteamId)
		//sqlQuery := fmt.Sprintf(`INSERT INTO game(name) VALUES ('%s')`, curName)
		_, err := DB.Exec(`INSERT INTO game(name) VALUES ($1)`, curName)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "game_name_key"` {
				fmt.Println("Error in InitGames", curSteamId, curName)
				return err
			}
		}
		//sqlQuery = fmt.Sprintf(`SELECT game.id FROM game WHERE name = '%s'`, curName)
		var gameid int

		err = DB.QueryRow(`SELECT game.id FROM game WHERE name = $1`, curName).Scan(&gameid)
		//fmt.Println(gameid)
		if err != nil {
			return err
		}
		//fmt.Println("Gameid:", gameid)
		//sqlQuery = fmt.Sprintf(`INSERT INTO gameprice(gameid, storeid, storegameid) VALUES (%d, %d, '%d')`, gameid, steamid, curSteamId)
		_, err = DB.Exec(`INSERT INTO gameprice(gameid, storeid, storegameid) VALUES ($1, $2, $3)`, gameid, steamid, curSteamId)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "gamepriceid"` {
				fmt.Println("Error in InitGames", curSteamId, curName)
				return err
			}
		}
	}
	return nil
}

func (DB *GameDB) InitGamePrice() error {
	gameids := []int{}
	steamgameids := []int{}
	sqlQuery := "SELECT storegameid, gameid FROM gameprice"
	//var temp *sql.DB = DB.PgDb
	rows, err := DB.Query(sqlQuery)
	if err != nil {
		return err
	}
	defer rows.Close()
	//fmt.Println("Parsing rows")
	// if !rows.Next() {
	// 	log.Fatalln("FML")
	// }
	// var res1 int
	// var res2 int

	// rows.Scan(&res1, &res2)

	for i := 0; rows.Next() && i < 250; i++ {
		//fmt.Printf("EZ")
		var steamgameid int
		var dbgameid int
		if err := rows.Scan(&steamgameid, &dbgameid); err != nil {
			return err
		}
		steamgameids = append(steamgameids, steamgameid)
		gameids = append(gameids, dbgameid)
	}
	//fmt.Println("SteamIDs: ", steamgameids)
	prices, err := steamapi.GetAppsPrice(&steamgameids, "ua")
	if err != nil {
		return err
	}
	//fmt.Println(len(*prices))
	for i := 0; i < len(*prices); i++ {
		sqlQuery = fmt.Sprintf(`UPDATE gameprice SET price = %f, discount = %f, free = %t WHERE gameid = %d AND storeid = %d`,
			(*(*prices)[i]).Initial/100,
			(*(*prices)[i]).Discount_percent,
			(*(*prices)[i]).Initial == 0,
			gameids[i],
			steamid)

		if (*prices)[i] == nil {
			continue
		}
		_, err = DB.Exec(`UPDATE gameprice SET price = $1, discount = $2, free = $3 WHERE gameid = $4 AND storeid = $5`,
			(*(*prices)[i]).Initial/100,
			(*(*prices)[i]).Discount_percent,
			(*(*prices)[i]).Initial == 0,
			gameids[i],
			steamid)
		if err != nil {
			return err
		}
	}

	return nil
}

type SolveDB struct {
	*sql.DB
}

func OpenSolve(credentials string) (*SolveDB, error) {
	db, err := sql.Open("postgres", credentials)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	for db.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			return nil, err
		}
	}
	fmt.Println("connected:", db.Ping() == nil)
	database := &SolveDB{db}
	return database, nil
}

func (Sol *SolveDB) SolveQuery() {
	var res1 int
	var res2 int

	rows, err := Sol.Query("SELECT storegameid, gameid FROM gameprice")
	if err != nil {
		log.Fatalln(err)
	}
	if !rows.Next() {
		log.Fatalln("Solve rows closed")
	}
	err = rows.Scan(&res1, &res2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Solve", res1, res2)
}

func GetAppPrice(gamename string, storename string) int {
	return 0
}

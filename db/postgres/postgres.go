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
	PgDb *sql.DB
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

func (DB GameDB) Close() {
	DB.PgDb.Close()
}

func (DB GameDB) InitTables() error {
	//needs to be reworked
	return nil
}

func (DB GameDB) InitGames() error {
	res := steamapi.GetAppList()
	for i := 0; i < 10; i++ {
		var curSteamId int = int(res[i].Get("appid").Value().(float64))
		var curName string = res[i].Get("name").Value().(string)
		//fmt.Println("Current game: ", curName, curSteamId)
		sqlQuery := fmt.Sprintf(`INSERT INTO game(name) VALUES ('%s')`, curName)
		_, err := DB.PgDb.Exec(sqlQuery)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "game_name_key"` {
				return err
			}
		}
		sqlQuery = fmt.Sprintf(`SELECT game.id FROM game WHERE name = '%s'`, curName)
		var gameid int

		err = DB.PgDb.QueryRow(sqlQuery).Scan(&gameid)
		fmt.Println(gameid)
		if err != nil {
			return err
		}
		//fmt.Println("Gameid:", gameid)
		sqlQuery = fmt.Sprintf(`INSERT INTO gameprice(gameid, storeid, storegameid) VALUES (%d, %d, '%d')`, gameid, steamid, curSteamId)
		_, err = DB.PgDb.Exec(sqlQuery)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "id"` {
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
	var temp *sql.DB = DB.PgDb
	rows, err := temp.Query(sqlQuery)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.NextResultSet() {
		fmt.Println(rows.Err())
	}
	fmt.Println("Parsing rows")
	for i := 0; rows.Next() && i < 250; i++ {
		var steamgameid int
		var dbgameid int
		if err := rows.Scan(&steamgameid, &dbgameid); err != nil {
			return err
		}
		steamgameids = append(steamgameids, steamgameid)
		gameids = append(gameids, dbgameid)
	}
	fmt.Println("SteamIDs: ", steamgameids)
	prices, err := steamapi.GetAppsPrice(&steamgameids, "ua")
	if err != nil {
		return err
	}

	for i := 0; i < len(*prices); i++ {
		sqlQuery = fmt.Sprintf(`INSERT INTO gameprice(gameid, storeid, price, discount, free)`, gameids[i], steamid,
			(*(*prices)[i]).Initial/100,
			(*(*prices)[i]).Discount_percent,
			(*(*prices)[i]).Initial == 0)
		_, err = DB.PgDb.Exec(sqlQuery)
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
	var database *SolveDB = &SolveDB{db}
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

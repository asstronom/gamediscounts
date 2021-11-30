package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gamediscounts/model/steamapi"
)

type Store int64

var (
	StoresNum int   = 1
	SteamID   Store = 1
)

type Country int64

var (
	UA Country = 0
	US Country = 1
)

func (c Country) CountryCode() string {
	switch c {
	case UA:
		return "ua"
	case US:
		return "us"
	}
	return "unknown"
}

func (c Country) Currency() string {
	switch c {
	case UA:
		return "UAH"
	case US:
		return "USD"
	}
	return "unknown"
}

type GameDB struct {
	*sql.DB
}

type GamePrice struct {
	gameid   int
	storeid  int
	initial  float64
	final    float64
	discount int
	isFree   bool
	currency string
}

type Game struct {
	name        string
	id          int
	price       []GamePrice
	description string
	imageURL    string
	genres      []string
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
	_, err = DB.Exec("DROP TABLE IF EXISTS featured CASCADE")
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

	_, err = DB.Exec(`CREATE TABLE featured (
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE)
		`)
	if err != nil {
		return err
	}
	return nil
}

func (DB *GameDB) InitStores() error {
	//query := fmt.Sprintf(`INSERT INTO store (id, name) VALUES (%d, 'steam')`, SteamID)
	_, err := DB.Exec(`INSERT INTO store (id, name) VALUES ($1, 'steam')`, SteamID)
	if err != nil {
		fmt.Println("Error InitStores")
		return err
	}

	return nil
}

func (DB *GameDB) InitGames() error {
	res := steamapi.GetAppList()
	for i := 0; i < len(res); i++ {
		var curSteamID int = int(res[i].Get("appid").Value().(float64))
		var curName string = res[i].Get("name").Value().(string)
		//fmt.Println("Current game: ", curName, curSteamID)
		//sqlQuery := fmt.Sprintf(`INSERT INTO game(name) VALUES ('%s')`, curName)
		_, err := DB.Exec(`INSERT INTO game(name) VALUES ($1)`, curName)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "game_name_key"` {
				fmt.Println("Error in InitGames", curSteamID, curName)
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
		//sqlQuery = fmt.Sprintf(`INSERT INTO gameprice(gameid, storeid, storegameid) VALUES (%d, %d, '%d')`, gameid, SteamID, curSteamID)
		_, err = DB.Exec(`INSERT INTO gameprice(gameid, storeid, storegameid) VALUES ($1, $2, $3)`, gameid, SteamID, curSteamID)
		if err != nil {
			if err.Error() != `pq: duplicate key value violates unique constraint "gamepriceid"` {
				fmt.Println("Error in InitGames", curSteamID, curName)
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

		if (*prices)[i] == nil {
			continue
		}
		_, err = DB.Exec(`UPDATE gameprice SET price = $1, discount = $2, free = $3, final = $6 WHERE gameid = $4 AND storeid = $5`,
			(*(*prices)[i]).Initial/100,
			(*(*prices)[i]).Discount_percent,
			(*(*prices)[i]).Initial == 0,
			gameids[i],
			SteamID,
			(*(*prices)[i]).Final/100)
		if err != nil {
			return err
		}
	}

	return nil
}

func (DB *GameDB) RefreshFeatured() error {
	ids, prices, err := steamapi.GetFeaturedCategories("ua")

	if err != nil {
		return err
	}
	_, err = DB.Exec(`TRUNCATE TABLE featured`)

	if err != nil {
		return err
	}

	for i := 0; i < len(prices); i++ {
		res := DB.QueryRow(`UPDATE gameprice SET price = $1, discount = $2, free = $3, final = $4 WHERE storegameid = $5 AND storeid = $6 RETURNING gameid`,
			prices[i].Initial/100,
			prices[i].Discount_percent,
			prices[i].Initial == 0,
			prices[i].Final,
			strconv.Itoa(ids[i]),
			SteamID)
		err = res.Err()
		if err != nil {
			return err
		}
		var curGameId int
		err = res.Scan(&curGameId)
		if err != nil {
			return err
		}
		_, err = DB.Exec(`INSERT INTO featured(gameid, storeid) VALUES ($1, $2)`, curGameId, SteamID)
		if err != nil {
			if err.Error() == `pq: duplicate key value violates unique constraint "featuredid"` {
				fmt.Printf("Duplicate game in featured Steam (gameid, storegameid): %d, %d\n", curGameId, ids[i])
			} else {
				return err
			}
		}
	}
	return nil
}

func (DB *GameDB) BestOffers(start int, count int, country Country) ([]GamePrice, error) {
	//cc := country.CountryCode()
	type featuredid struct {
		gameid  int
		storeid int
	}
	var (
		ids []featuredid
		res []GamePrice
	)
	rows, err := DB.Query(`SELECT gameid, storeid FROM featured LIMIT $1 OFFSET $2`, count, start)
	if err != nil {
		return nil, err
	}
	for i := 0; rows.Next(); i++ {
		var curid int
		var curstore int
		rows.Scan(&curid, &curstore)
		ids = append(ids, featuredid{gameid: curid, storeid: curstore})
	}

	for i := 0; i < len(ids); i++ {
		row := DB.QueryRow(`
		SELECT gameid, storeid, price, final, discount, free
		FROM gameprice 
		WHERE
		gameid = $1 AND storeid = $2`,
			ids[i].gameid, ids[i].storeid)
		if row.Err() != nil {
			return nil, err
		}
		temp := GamePrice{}
		err = row.Scan(&temp.gameid, &temp.storeid, &temp.initial, &temp.final, &temp.discount, &temp.isFree)
		if err != nil {
			return nil, err
		}
		temp.currency = country.Currency()
		res = append(res, temp)
	}
	return res, nil
}

// type SolveDB struct {
// 	*sql.DB
// }

// func OpenSolve(credentials string) (*SolveDB, error) {
// 	db, err := sql.Open("postgres", credentials)
// 	if err != nil {
// 		return nil, err
// 	}

// 	start := time.Now()
// 	for db.Ping() != nil {
// 		if start.After(start.Add(10 * time.Second)) {
// 			return nil, err
// 		}
// 	}
// 	fmt.Println("connected:", db.Ping() == nil)
// 	database := &SolveDB{db}
// 	return database, nil
// }

// func (Sol *SolveDB) SolveQuery() {
// 	var res1 int
// 	var res2 int

// 	rows, err := Sol.Query("SELECT storegameid, gameid FROM gameprice")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	if !rows.Next() {
// 		log.Fatalln("Solve rows closed")
// 	}
// 	err = rows.Scan(&res1, &res2)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Println("Solve", res1, res2)
// }

func (DB *GameDB) GetAppPrice(gameid int, storeid int, country Country) (GamePrice, error) {
	var res GamePrice
	row := DB.QueryRow(`SELECT gameid, storeid, price, final, discount, free FROM gameprice WHERE gameid = $1 AND storeid = $2`, gameid, storeid)
	if row.Err() != nil {
		return GamePrice{}, row.Err()
	}
	err := row.Scan(&res.gameid, &res.storeid, &res.initial, &res.final, &res.discount, &res.isFree)
	res.currency = country.Currency()
	if err != nil {
		return GamePrice{}, err
	}
	return res, nil
}

func (DB *GameDB) GetGameName(gameid int) (string, error) {
	var res string
	row := DB.QueryRow(`SELECT name FROM game WHERE id = $1`, gameid)
	if row.Err() != nil {
		return "", row.Err()
	}
	err := row.Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (DB *GameDB) GetGame(gameid int, country Country) (Game, error) {
	var res Game
	row := DB.QueryRow(`SELECT name, id FROM game WHERE id = $1`, gameid)
	if row.Err() != nil {
		return Game{}, row.Err()
	}
	err := row.Scan(&res.name, &res.id)
	if err != nil {
		return Game{}, err
	}
	for i := 1; i <= StoresNum; i++ {
		temp, err := DB.GetAppPrice(gameid, i, country)
		if err != nil {
			//return Game{}, err
		}
		res.price = append(res.price, temp)
	}
	return res, nil
}

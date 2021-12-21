package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gamediscounts/model/steamapi"
	_ "github.com/lib/pq"
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
	Gameid   int     `json:"gameid"`
	Storeid  int     `json:"storeid"`
	Initial  float64 `json:"initial"`
	Final    float64 `json:"final"`
	Discount int     `json:"discount"`
	IsFree   bool    `json:"is_free"`
	Currency string  `json:"currency"`
}

type DLC struct {
	Name        string
	Id          int
	GameID      int
	Description string
	HeaderImage string
}

type Package struct {
	Name      string
	Id        int
	SmallLogo string
	PageImage string
	GameIDs   []int
	DLCIDs    []int
}

type Game struct {
	Name        string      `json:"name"`
	Id          int         `json:"id"`
	Price       []GamePrice `json:"price"`
	Description string      `json:"description"`
	ImageURL    string      `json:"image_url"`
	Genres      []string    `json:"genres"`
	Packages    []int
	DLCs        []int
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
	_, err = DB.Exec("DROP TABLE IF EXISTS dlc CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS dlcprice CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS package CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS packageprice CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS featured CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS genre CASCADE")
	if err != nil {
		return err
	}
	_, err = DB.Exec("DROP TABLE IF EXISTS gamegenre CASCADE")
	if err != nil {
		return err
	}

	_, err = DB.Exec("DROP TABLE IF EXISTS gamepackage CASCADE")
	if err != nil {
		return err
	}

	_, err = DB.Exec("DROP TABLE IF EXISTS dlcpackage CASCADE")
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE game (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE,
		description TEXT,
		headerimage TEXT DEFAULT 'https:\/\/cdn.akamai.steamstatic.com\/steam\/apps\/292030\/header_russian.jpg?t=1621939214')`)
	if err != nil {
		fmt.Println("creating game")
		return err
	}
	_, err = DB.Exec(`CREATE TABLE store (
		id SERIAL PRIMARY KEY,
		name text UNIQUE,
		icon text,
		banner text)`)
	if err != nil {
		fmt.Println("creating store")
		return err
	}
	_, err = DB.Exec(`CREATE TABLE gameprice (
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE,
		CONSTRAINT gamePriceId PRIMARY KEY (gameid, storeid),
		storegameid text UNIQUE,
		initial NUMERIC,
		final NUMERIC, 
		discount INT DEFAULT 0,
		free BOOLEAN DEFAULT FALSE,
		currency VARCHAR (10))`)
	if err != nil {
		fmt.Println("creating gameprice")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE featured (
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE)
		`)
	if err != nil {
		fmt.Println("creating featured")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE genre (
		id SERIAL PRIMARY KEY,
		name text UNIQUE)`)
	if err != nil {
		fmt.Println("creating genre")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE gamegenre (
		genreid INT REFERENCES genre (id) ON UPDATE CASCADE ON DELETE CASCADE,
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		PRIMARY KEY (genreid, gameid))`)
	if err != nil {
		fmt.Println("creating gamegenre")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE dlc (
		id SERIAL PRIMARY KEY,
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		name text UNIQUE,
		description TEXT,
		headerImage TEXT DEFAULT 'https:\/\/cdn.akamai.steamstatic.com\/steam\/apps\/292030\/header_russian.jpg?t=1621939214')`)
	if err != nil {
		fmt.Println("creating dlc")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE dlcprice (
		dlcid INT REFERENCES dlc (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE,
		CONSTRAINT dlcPriceId PRIMARY KEY (dlcid, storeid),
		storedlcid text UNIQUE,
		initial NUMERIC,
		final NUMERIC, 
		discount INT DEFAULT 0,
		free BOOLEAN DEFAULT FALSE,
		currency VARCHAR (10))`)
	if err != nil {
		fmt.Println("creating dlcprice")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE package (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE,
		smalllogo TEXT,
		pageimage TEXT)`)
	if err != nil {
		fmt.Println("creating package")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE gamepackage (
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		packageid INT REFERENCES package (id) ON UPDATE CASCADE ON DELETE CASCADE)`)
	if err != nil {
		fmt.Println("creating gamepackage")
		return err
	}

	_, err = DB.Exec(`CREATE TABLE dlcpackage (
		dlcid INT REFERENCES dlc (id) ON UPDATE CASCADE ON DELETE CASCADE,
		packageid INT REFERENCES package (id) ON UPDATE CASCADE ON DELETE CASCADE)`)
	if err != nil {
		fmt.Println("creating dlcpackage")
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

	sqlQuery := "SELECT storegameid, gameid FROM gameprice"
	//var temp *sql.DB = DB.PgDb
	rows, err := DB.Query(sqlQuery)
	if err != nil {
		return err
	}
	sqlQuery = "SELECT COUNT(*) FROM game"
	counterRow := DB.QueryRow(sqlQuery)
	gamecounter := 0
	counterRow.Scan(&gamecounter)
	defer rows.Close()
	progressCounter := 0
	for {
		gameids := []int{}
		steamgameids := []int{}
		fmt.Println(progressCounter)
		for i := 0; rows.Next() && i < 250; i++ {
			var steamgameid int
			var dbgameid int
			if err := rows.Scan(&steamgameid, &dbgameid); err != nil {
				return err
			}
			steamgameids = append(steamgameids, steamgameid)
			gameids = append(gameids, dbgameid)
		}
		if (len(steamgameids)) == 0 {
			break
		}
		prices, err := steamapi.GetAppsPrice(&steamgameids, "ua")
		if err != nil {
			return err
		}
		for i := 0; i < len(*prices); i++ {
			progressCounter++
			if len(*prices) == 0 {
				break
			}

			if (*prices)[i] == nil {
				continue
			}
			_, err = DB.Exec(`UPDATE gameprice SET price = $1, discount = $2, free = $3, final = $6, correct = $7 WHERE gameid = $4 AND storeid = $5`,
				(*(*prices)[i]).Initial/100,
				(*(*prices)[i]).Discount_percent,
				false,
				gameids[i],
				SteamID,
				(*(*prices)[i]).Final/100,
				true)

			if err != nil {
				return err
			}

		}
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func (DB *GameDB) BestOffers(start int, count int, country Country) ([]GamePrice, error) {
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
		SELECT gameid, storeid, initial, final, discount, free
		FROM gameprice 
		WHERE
		gameid = $1 AND storeid = $2`,
			ids[i].gameid, ids[i].storeid)
		if row.Err() != nil {
			return nil, err
		}
		temp := GamePrice{}
		err = row.Scan(&temp.Gameid, &temp.Storeid, &temp.Initial, &temp.Final, &temp.Discount, &temp.IsFree)
		if err != nil {
			return nil, err
		}
		temp.Currency = country.Currency()
		res = append(res, temp)
	}
	return res, nil
}

func (DB *GameDB) GetAppPrice(gameid int, storeid int, country Country) (GamePrice, error) {
	var res GamePrice
	row := DB.QueryRow(`SELECT gameid, storeid, price, final, discount, free FROM gameprice WHERE gameid = $1 AND storeid = $2`, gameid, storeid)
	if row.Err() != nil {
		return GamePrice{}, row.Err()
	}
	err := row.Scan(&res.Gameid, &res.Storeid, &res.Initial, &res.Final, &res.Discount, &res.IsFree)
	res.Currency = country.Currency()
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
	err := row.Scan(&res.Name, &res.Id)
	if err != nil {
		return Game{}, err
	}
	for i := 1; i <= StoresNum; i++ {
		temp, err := DB.GetAppPrice(gameid, i, country)
		if err != nil {
			//return Game{}, err
		}
		res.Price = append(res.Price, temp)
	}
	return res, nil
}

func getSteamIDs() []int {
	resultids := steamapi.GetAppList()
	steamgameids := []int{}
	for _, r := range resultids {
		steamgameids = append(steamgameids, int(r.Get("appid").Int()))
	}
	return steamgameids
}

func (DB *GameDB) insertPackage(pack steamapi.Package) error {
	var packageid int
	row := DB.QueryRow("INSERT INTO package (name, smalllogo, pageimage) VALUES ($1, $2, $3) RETURNING id",
		pack.Name, pack.SmallLogo, pack.PageImage)
	err := row.Err()
	if err != nil {
		return err
	}
	row.Scan(&packageid)
	_, err = DB.Exec(`INSERT INTO packageprice (packageid, storeid, storepackageid, initial, final, discount_percent, currency, isfree, individual)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		packageid, SteamID, pack.PackageId, pack.Price.Initial, pack.Price.Final, pack.Price.Discount_percent, pack.Price.Currency, pack.Price.IsFree, pack.Individual)
	if err != nil {
		return err
	}
	for _, v := range pack.AppIds {
		app, err := steamapi.GetAppInfo(v, "ua")
		if err != nil {
			continue
		}
		switch app.Type {
		case steamapi.Game:
			var gameid int
			row := DB.QueryRow(`SELECT id FROM game WHERE name = $1`, app.Name)
			if row.Err() != nil {
				return row.Err()
			}
			row.Scan(&gameid)
		default:
			var dlcid int
			row := DB.QueryRow(`SELECT id FROM dlc WHERE name = $1`, app.Name)
			if row.Err() != nil {
				return row.Err()
			}
			row.Scan(&dlcid)
			log.Println("dlcid: ", dlcid)
		}
	}

	return nil
}

func (DB *GameDB) insertDLC(app steamapi.AppInfo, gameid int) error {
	if app.Type == steamapi.Game {
		return fmt.Errorf("game not a dlc")
	}
	var row *sql.Row
	if gameid != -1 {
		row = DB.QueryRow("INSERT INTO dlc (gameid, name, description, headerimage) VALUES ($4, $1, $2, $3) RETURNING id",
			app.Name, app.Description, app.HeaderImage, gameid,
		)
	} else {
		row = DB.QueryRow("INSERT INTO dlc (name, description, headerimage) VALUES ($1, $2, $3) RETURNING id",
			app.Name, app.Description, app.HeaderImage,
		)
	}
	if row.Err() != nil {
		if row.Err().Error() == `pq: duplicate key value violates unique constraint "dlc_name_key"` {
			row = DB.QueryRow(`SELECT id FROM dlc WHERE name = $1`, app.Name)
			if row.Err() != nil {
				return row.Err()
			}
		} else {
			return row.Err()
		}
	}
	var dlcid int
	err := row.Scan(&dlcid)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`INSERT INTO dlcprice (dlcid, storeid, storedlcid, initial, final, discount, currency, free)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		dlcid, SteamID, app.Appid, app.Price.Initial, app.Price.Final, app.Price.Discount_percent, app.Price.Currency, app.Price.IsFree,
	)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "dlcpriceid"` {
			_, err = DB.Exec(`UPDATE dlcprice SET initial = $1, final = $2, discount = $3, currency = $4, free = $5 WHERE dlcid = $6 AND storeid = $7`,
				app.Price.Initial, app.Price.Final, app.Price.Discount_percent, app.Price.Currency, app.Price.IsFree, dlcid, SteamID,
			)
			if err != nil {
				fmt.Println("error updating existing dlcprice")
				return err
			}
		} else {
			log.Println("dlc id: ", dlcid)
			return err
		}

	}
	return nil
}

func (DB *GameDB) linkGameGenre(genreid, gameid int) error {
	_, err := DB.Exec(`INSERT INTO gamegenre (gameid, genreid) VALUES ($1, $2)`, gameid, genreid)
	if err != nil {
		return err
	}
	return nil
}

func convertSteamPrice(price steamapi.PriceOverview, gameid int, storeid int) GamePrice {
	var result GamePrice
	result.Gameid = gameid
	result.Storeid = storeid
	result.Initial = price.Initial
	result.Final = price.Final
	result.Discount = int(price.Discount_percent)
	result.IsFree = price.IsFree
	result.Currency = price.Currency
	return result
}

func (DB *GameDB) matchGenre(genre string) (int, error) {
	var result int
	row := DB.QueryRow(`SELECT id FROM genre WHERE name = $1`, genre)
	if row.Err() != nil {
		return result, row.Err()
	}
	err := row.Scan(&result)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			result = 0
		} else {
			return result, err
		}
	}
	if result == 0 {
		row := DB.QueryRow(`INSERT INTO genre (name) VALUES ($1) RETURNING id`, genre)
		if row.Err() != nil {
			return result, row.Err()
		}
		err := row.Scan(&result)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (DB *GameDB) insertGamePrice(price GamePrice, storegameid string) error {
	_, err := DB.Exec(`INSERT INTO gameprice (gameid, storeid, storegameid, initial, final, discount, free, currency) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		price.Gameid, price.Storeid, storegameid, price.Initial, price.Final, price.Discount, price.IsFree, price.Currency,
	)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "gamepriceid"` {
			fmt.Println("duplicate gameprice")
			_, err = DB.Exec(`UPDATE gameprice SET initial = $1, final = $2, discount = $3, free = $4, currency = $5 WHERE gameid = $6 AND storeid = $7`,
				price.Initial, price.Final, price.Discount, price.IsFree, price.Currency, price.Gameid, price.Storeid,
			)
			if err != nil {
				fmt.Println("error updating gameprice")
				return err
			}
		} else {
			fmt.Println("other error while inserting gameprice")
			return err
		}
	}
	return nil
}

func (DB *GameDB) insertGame(app steamapi.AppInfo) error {
	if app.Type != steamapi.Game {
		return fmt.Errorf("not a game")
	}
	var appid int
	row := DB.QueryRow(`INSERT INTO game (name, description, headerimage) VALUES ($1, $2, $3) RETURNING id`, app.Name, app.Description, app.HeaderImage)
	if row.Err() != nil {
		if row.Err().Error() == `pq: duplicate key value violates unique constraint "game_name_key"` {
			row = DB.QueryRow(`SELECT id FROM game WHERE name = $1`, app.Name)
		} else {
			return row.Err()
		}
	}
	err := row.Scan(&appid)
	if err != nil {
		return err
	}
	for _, v := range app.Genres {
		genreid, err := DB.matchGenre(v)
		if err != nil {
			log.Println("error in matching genre")
			return err
		}
		DB.linkGameGenre(genreid, appid)
	}
	price := convertSteamPrice(app.Price, appid, int(SteamID))
	err = DB.insertGamePrice(price, strconv.Itoa(app.Appid))
	if err != nil {
		log.Println("error in inserting gameprice")
		return err
	}

	for _, v := range app.DLC {
		dlcinfo, err := steamapi.GetAppInfo(v, "ua")
		if err != nil {
			log.Println("error in getting dlcinfo")
			continue
		}
		err = DB.insertDLC(dlcinfo, appid)
		if err != nil {
			log.Println("error in inserting dlc")
			return err
		}
	}

	return nil
}

func (DB *GameDB) InitDatabase() error {
	var counter int
	steamgameids := getSteamIDs()
	for i, v := range steamgameids {

		app, err := steamapi.GetAppInfo(v, "ua")
		if err != nil {
			fmt.Println(err, v, i, "/", len(steamgameids))
			continue
		}
		counter++
		fmt.Println("Valid app counter: ", counter, i, "/", len(steamgameids))
		switch app.Type {
		case steamapi.Game:
			err := DB.insertGame(app)
			if err != nil {
				fmt.Println("error inserting game")
				return err
			}
		default:
			err := DB.insertDLC(app, -1)
			if err != nil {
				fmt.Println("error inserting dlc without game")
			}
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
	loopstart:
		var res *sql.Row
		res = DB.QueryRow(`UPDATE gameprice SET initial = $1, discount = $2, free = $3, final = $4 WHERE storegameid = $5 AND storeid = $6 RETURNING gameid`,
			prices[i].Initial,
			prices[i].Discount_percent,
			prices[i].IsFree,
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
			if err.Error() == `sql: no rows in result set` {
				appinfo, err := steamapi.GetAppInfo(ids[i], "ua")
				if err != nil {
					continue
				}
				err = DB.insertGame(appinfo)
				if err != nil {
					fmt.Println("error inserting game in refresh featured")
					return err
				}
				goto loopstart

			} else {
				fmt.Println("error while scanning featured")
				return err
			}
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

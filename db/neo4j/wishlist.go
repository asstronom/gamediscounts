package wishlist

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type WishlistDB struct {
	db neo4j.Driver
}

func OpenDB(uri string, username string, password string) (*WishlistDB, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return &WishlistDB{db: driver}, nil
}

func (DB *WishlistDB) Clear() error {
	session := DB.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	_, err := session.Run("MATCH ()-[t:tracks]-() DELETE t", map[string]interface{}{})
	if err != nil {
		return err
	}
	_, err = session.Run("MATCH (u:User), (g:Game) DELETE u, g", map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (DB *WishlistDB) AddUser(username string) error {
	session := DB.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	_, err := session.Run("CREATE (u:User{username: $username})", map[string]interface{}{"username": username})
	if err != nil {
		if neo4j.IsNeo4jError(err) {
			if (err.(*neo4j.Neo4jError).Code) == "Neo.ClientError.Schema.ConstraintValidationFailed" {
				return nil
			}
		}
		return err
	}
	return nil
}

func (DB *WishlistDB) AddGame(gameid int) error {
	session := DB.db.NewSession((neo4j.SessionConfig{}))
	defer session.Close()
	_, err := session.Run("CREATE (g:Game{id: $gameid})", map[string]interface{}{"gameid": gameid})
	if err != nil {
		if neo4j.IsNeo4jError(err) {
			if (err.(*neo4j.Neo4jError).Code) == "Neo.ClientError.Schema.ConstraintValidationFailed" {
				return nil
			}
		}
		return err
	}
	return nil
}

func (DB *WishlistDB) CheckIfUserExists(username string) error {
	session := DB.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.Run(`MATCH (u:User{username:$username}) RETURN u`, map[string]interface{}{"username": username})
	if err != nil {
		return err
	}
	records, err := result.Collect()
	if len(records) == 0 {
		return fmt.Errorf("User doesn't exist")
	}
	return nil
}

func (DB *WishlistDB) CheckIfGameExists(gameid int) error {
	return nil
}

func (DB *WishlistDB) GetWishlist(username string) ([]int64, error) {
	res := []int64{}
	session := DB.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	err := DB.CheckIfUserExists(username)

	if err != nil {
		return nil, err
	}

	records, err := session.Run(`MATCH (u:User{username:$username})-->(n) RETURN n`, map[string]interface{}{"username": username})
	if err != nil {
		return nil, err
	}
	result, err := records.Collect()
	if err != nil {
		return nil, err
	}
	for _, r := range result {
		res = append(res, r.Values[0].(neo4j.Node).Props["id"].(int64))
	}
	return res, nil
}

func (DB *WishlistDB) AddGameToWishList(username string, gameid int) error {
	session := DB.db.NewSession(neo4j.SessionConfig{})

	defer session.Close()
	_, err := session.Run(`MERGE (u:User{username: $username}) MERGE (g:Game{id: $gameid})
	 MERGE (u)-[t:tracks]->(g)
	 ON MATCH
	 SET t.id = 0`,
		map[string]interface{}{"username": username, "gameid": gameid})
	if err != nil {
		return err
	}
	return nil
}

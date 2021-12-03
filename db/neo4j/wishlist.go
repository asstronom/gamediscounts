package wishlist

import (
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
		return err
	}
	return nil
}

func (DB *WishlistDB) AddGame(gameid int) error {
	session := DB.db.NewSession((neo4j.SessionConfig{}))
	defer session.Close()

	_, err := session.Run("CREATE (g:Game{id: $gameid})", map[string]interface{}{"gameid": gameid})
	if err != nil {
		return err
	}
	return nil
}

func (DB *WishlistDB) AddGameToWishList(username string, gameid int) error {
	session := DB.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	_, err := session.Run("MATCH (u:User{username: $username}), (g:Game{id: $gameid}) CREATE (u)-[t:tracks]->(g)",
		map[string]interface{}{"username": username, "gameid": gameid})
	if err != nil {
		return err
	}
	return nil
}

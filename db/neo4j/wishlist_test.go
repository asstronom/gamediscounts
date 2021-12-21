package wishlist

import (
	"log"
	"testing"
)

const (
	wishlistURI  = "neo4j://localhost:7687"
	wishUsername = "neo4j"
	wishPassword = "GuesgP4LPLS"
)

func TestRemoveSingleTrack(t *testing.T) {
	wishlistDB, er := OpenDB(wishlistURI, wishUsername, wishPassword)
	if er != nil {
		log.Fatalln("error opening wishlist: ", er)
	}
	wishlistDB.AddGameToWishList("asstronom", 1)
	wishlistDB.AddGameToWishList("pudgebooster", 1)
	wishlistDB.AddGameToWishList("asstronom", 2)
	wishlistDB.AddGameToWishList("pudgebooster", 2)
	wishlistDB.AddGameToWishList("asstronom", 3)
	wishlistDB.RemoveSingleTrack("asstronom", 1)
	wishlistDB.RemoveWholeWishlist("asstronom")
}

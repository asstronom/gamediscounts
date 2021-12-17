package server

import (
	"encoding/json"
	"fmt"
	"github.com/gamediscounts/auth"
	"github.com/gamediscounts/db/postgres"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) HandleIndex() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		fmt.Println(values)
		start, err := strconv.Atoi(values.Get("start"))
		if err != nil {
			log.Println(err)
		}
		count, err := strconv.Atoi(values.Get("count"))
		if err != nil {
			log.Println(err)
		}
		offers, err := s.gameDB.BestOffers(start, count, postgres.UA)
		fmt.Println(len(offers)) //check len of array for debugging
		if err != nil {
			log.Println(err)
		}
		s.respond(w, r, offers, http.StatusOK)
	}
}
func (s *Server) HandleSingleGame() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		username, err := auth.GetTokenUsername(r)
		idint, _ := strconv.Atoi(id)
		game, err := s.gameDB.GetGame(idint, postgres.UA)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(vars) // just for debug
		fmt.Println(username)
		s.respond(w, r, game, http.StatusOK)
	}
}
func (s *Server) WishlistAddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		gameID, err := strconv.Atoi(id)
		fmt.Println("game id to add:", gameID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		username, err := auth.GetTokenUsername(r)
		err = s.wishDB.AddGameToWishList(username, gameID)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
func (s *Server) WishlistAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username, err := auth.GetTokenUsername(r)
		//wishlistDB, err := wishlist.OpenDB(wishlist.WishlistURI, wishlist.WishUsername, wishlist.WishPassword)
		wishlistIDSlice, err := s.wishDB.GetWishlist(username)
		fmt.Println(wishlistIDSlice) // for debug
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var games []postgres.Game
		for _, gameID := range wishlistIDSlice {
			game, err := s.gameDB.GetGame(int(gameID), postgres.UA)
			if err != nil {
				log.Fatalln(err)
			}
			games = append(games, game)
		}
		fmt.Println(games)
		err = json.NewEncoder(w).Encode(games)
	}
}

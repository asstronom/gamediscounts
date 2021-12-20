package server

import (
	"context"
	"encoding/json"
	userdb "github.com/gamediscounts/db/couchdb"
	wishlist "github.com/gamediscounts/db/neo4j"
	"github.com/gamediscounts/db/postgres"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router  *mux.Router
	context context.Context
	gameDB  *postgres.GameDB
	userDB  *userdb.UserDB
	wishDB  *wishlist.WishlistDB
}

func Init(ctx context.Context, gameDB *postgres.GameDB, userDB *userdb.UserDB, wishDB *wishlist.WishlistDB) *Server {
	router := mux.NewRouter()
	s := &Server{
		router,
		ctx,
		gameDB,
		userDB,
		wishDB,
	}
	s.routes()
	return s
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		//dataJSON, _ := json.Marshal(data)
		//lol := string(dataJSON) // lol:)
		//fmt.Println(string(dataJSON)) // for debugging
		err := json.NewEncoder(w).Encode(data)
		if err != nil {

		}
	}
}
func (s *Server) error(w http.ResponseWriter, r *http.Request, err error, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err != nil {
		err = json.NewEncoder(w).Encode(e(err))
		if err != nil {
			//
		}
	}
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

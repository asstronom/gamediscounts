package server

import (
	"github.com/gamediscounts/db/postgres"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) HandleIndex() http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request){
		values := r.URL.Query()
		start, err := strconv.Atoi(values.Get("start"))
		if err != nil {
			log.Println(err)
		}
		count, err := strconv.Atoi(values.Get("start"))
		if err != nil {
			log.Println(err)
		}
		offers, err := s.db.BestOffers(start,count,postgres.UA)
		if err != nil{
			log.Println(err)
		}
		s.respond(w,r,offers,http.StatusOK)
	}
}
func (s *Server) HandleSingleGame() http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id:= vars["id"]
		idint, _  := strconv.Atoi(id)
		name, err := s.db.GetGameName(idint)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(vars) // just for debug
		s.respond(w, r,map[string]interface{}{
			"app":name,
		},http.StatusOK)
	}
}


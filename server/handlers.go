package server

import (
	"fmt"
	"github.com/gamediscounts/db/postgres"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) HandleIndex() http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request){
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
		offers, err := s.db.BestOffers(start,count,postgres.UA)
		fmt.Println(len(offers)) //check len of array for debugging
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


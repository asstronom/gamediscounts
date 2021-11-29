package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) HandleIndex() http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request){
		s.respond(w,r,map[string]interface{}{
			"message":"hello",
		},http.StatusOK)
	}
}
func (s *Server) HandleSingleGame() http.HandlerFunc{

	return func (w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id:= vars["id"]
		idint, _  := strconv.Atoi(id)
		name, err := s.db.GetGameName(idint)
		if err != nil {
			log.Fatalln(err)	
		}
		//fmt.Println(vars) // just for debug
		s.respond(w, r,map[string]interface{}{
			"app":name,
		},http.StatusOK)
	}
}


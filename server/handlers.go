package server

import (
	"github.com/gorilla/mux"
	"net/http"
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
		//fmt.Println(vars) // just for debug
		s.respond(w, r,map[string]interface{}{
			"app":id,
		},http.StatusOK)
	}
}


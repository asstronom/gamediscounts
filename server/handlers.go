package server

import "net/http"

func (s *Server) HandleIndex() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		s.respond(w,r,map[string]interface{}{
			"message":"hello",
		},200)
	}
}


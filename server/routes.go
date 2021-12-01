package server

func (s *Server) routes()  {
	s.router.HandleFunc("/",s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/featured",s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/app/{id}",s.HandleSingleGame()).Methods("GET")
}

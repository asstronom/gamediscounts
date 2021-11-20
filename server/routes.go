package server

func (s *Server) routes()  {
	s.router.HandleFunc("/",s.HandleIndex())
	s.router.HandleFunc("/featured",s.HandleIndex())

}

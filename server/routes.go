package server

import "github.com/gamediscounts/auth"

func (s *Server) routes()  {
	s.router.HandleFunc("/",s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/featured",s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/app/{id}",s.HandleSingleGame()).Methods("GET")

	//auth
	s.router.HandleFunc("/register",auth.SignUp).Methods("POST")
	s.router.HandleFunc("/login",auth.SignIn).Methods("POST")
	//  /wishilist/:id PUT
	//  /wishilist/:


}

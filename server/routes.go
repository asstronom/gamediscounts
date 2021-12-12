package server

import "github.com/gamediscounts/auth"

func (s *Server) routes() {
	s.router.Handle("/", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/featured", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/app/{id}", s.HandleSingleGame()).Methods("GET")
	//wishlish
	s.router.HandleFunc("/wishlist/{id}", auth.IsAuthorized()).Methods("PUT")
	s.router.HandleFunc("/wishlist/", auth.IsAuthorized()).Methods("GET")
	
	//auth
	s.router.HandleFunc("/register", auth.SignUp).Methods("POST")
	s.router.HandleFunc("/login", auth.SignIn).Methods("POST")
	//  /wishilist/:id PUT
	//  /wishilist/:
}

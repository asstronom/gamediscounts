package server

import "github.com/gamediscounts/auth"

func (s *Server) routes() {
	s.router.Handle("/", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/featured", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/app/{id}", auth.IsAuthorized(s.HandleSingleGame())).Methods("GET")
	//wishlish
	s.router.HandleFunc("/wishlist/{id}", auth.IsAuthorized(s.WishlistAddItem())).Methods("PUT")
	//s.router.HandleFunc("/wishlist/{id}", auth.IsAuthorized()).Methods("DELETE")
	s.router.HandleFunc("/wishlist", auth.IsAuthorized(s.WishlistAll())).Methods("GET")
	//s.router.HandleFunc("/wishlist", auth.IsAuthorized()).Methods("DELETE")

	//auth
	s.router.HandleFunc("/register", auth.SignUp).Methods("POST")
	s.router.HandleFunc("/login", auth.SignIn).Methods("POST")
	s.router.HandleFunc("/refresh", auth.Refresh).Methods("GET")
}

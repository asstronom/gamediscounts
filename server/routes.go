package server

import "github.com/gamediscounts/auth"

func (s *Server) routes() {
	s.router.Handle("/", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/featured", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/app/{id}", s.HandleSingleGame()).Methods("GET")
	//wishlish
	s.router.HandleFunc("/wishlist/{id}", auth.IsAuthorized(s.WishlistAddItem())).Methods("PUT")
	s.router.HandleFunc("/wishlist/{id}", auth.IsAuthorized(s.WishlistRemoveItem())).Methods("DELETE")
	s.router.HandleFunc("/wishlist", auth.IsAuthorized(s.WishlistAll())).Methods("GET")
	s.router.HandleFunc("/username", auth.IsAuthorized(auth.FetchUserName)).Methods("GET")
	//s.router.HandleFunc("/wishlist", auth.IsAuthorized()).Methods("DELETE")

	//auth
	s.router.HandleFunc("/register", auth.SignUp).Methods("POST")
	s.router.HandleFunc("/logout", auth.Logout).Methods("POST")
	s.router.HandleFunc("/login", auth.SignIn).Methods("POST")
	s.router.HandleFunc("/refresh", auth.Refresh).Methods("GET")
	//wishlish notification
	s.router.HandleFunc("/notify", s.Notify()).Methods("GET")

}

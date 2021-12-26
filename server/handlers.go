package server

import (
	"encoding/json"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gamediscounts/auth"
	"github.com/gamediscounts/db/postgres"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func (s *Server) HandleIndex() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
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

		offers, err := s.gameDB.BestOffers(start, count, postgres.UA)
		// var games []postgres.Game
		// for _, item := range offers {
		// 	game, err := s.gameDB.GetGame(item.Gameid, postgres.UA)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	games = append(games, game)
		//}
		fmt.Println(len(offers)) //check len of array for debugging
		//fmt.Println(games)       //check len of array for debugging
		if err != nil {
			log.Println(err)
		}
		s.respond(w, r, offers, http.StatusOK)
	}
}
func (s *Server) HandleSingleGame() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		//username, err := auth.GetTokenUsername(r)
		idint, _ := strconv.Atoi(id)
		game, err := s.gameDB.GetGame(idint, postgres.UA)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(vars) // just for debug
		//fmt.Println(username)
		s.respond(w, r, game, http.StatusOK)
	}
}
func (s *Server) WishlistAddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]
		gameID, err := strconv.Atoi(id)
		fmt.Println("game id to add:", gameID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		username, err := auth.GetTokenUsername(r)
		err = s.wishDB.AddGameToWishList(username, gameID)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
func (s *Server) WishlistRemoveItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		gameID, err := strconv.Atoi(id)
		fmt.Println("id game of game to be removed:", gameID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		username, err := auth.GetTokenUsername(r)
		err = s.wishDB.RemoveSingleTrack(username, gameID)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
func (s *Server) WishlistAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Method", "true")

		username, err := auth.GetTokenUsername(r)
		//wishlistDB, err := wishlist.OpenDB(wishlist.WishlistURI, wishlist.WishUsername, wishlist.WishPassword)
		wishlistIDSlice, err := s.wishDB.GetWishlist(username)
		fmt.Println(wishlistIDSlice) // for debug
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var games []postgres.Game
		for _, gameID := range wishlistIDSlice {
			game, err := s.gameDB.GetGame(int(gameID), postgres.UA)
			if err != nil {
				log.Fatalln(err)
			}
			games = append(games, game)
		}
		fmt.Println(games)
		err = json.NewEncoder(w).Encode(games)
	}
}
func (s *Server) Notify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := s.wishDB.GetGames()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(games)

		var discounted []int
		for _, game := range games {
			item, err2 := s.gameDB.GetGame(int(game), postgres.UA)
			if err2 != nil {
				log.Println(err2)
			}
			if item.Price[0].Discount > 0 {
				discounted = append(discounted, int(game))
			}
		}
		fmt.Println(discounted)
		for _, item := range discounted {
			userToBeNotified, err := s.wishDB.GetUsersByGame(item)
			if err != nil {
				log.Println(err)
			}
			var Emails []string
			for _, username := range userToBeNotified {
				user, err := s.userDB.GetUserByName(username)
				if err != nil {
					log.Println(err)
				}
				Emails = append(Emails, fmt.Sprintf(user.EmailName+"@"+user.EmailDomain))
			}
			fmt.Println(Emails)
			game, err := s.gameDB.GetGame(item, postgres.UA)
			if err != nil {
				log.Println(err)
			}
			err = SendEmailNotification(Emails, game)
			if err != nil {
				log.Println(err)
			}
		}

	}
}
func SendEmailNotification(emailStrSlice []string, game postgres.Game) error {
	//	fmt.Println(emailStrSlice)
	server := mail.SMTPServer{}
	server.KeepAlive = true
	server.Host = "smtp.gmail.com"
	server.Port = 587
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	server.Username = os.Getenv("MAIL")
	server.Password = os.Getenv("MAIL_PASS")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range emailStrSlice {

		email := mail.NewMSG()
		email.SetFrom("Gamediscounts <gamediscountsiasa@gmail.com>")
		email.AddTo(item)

		info, err := json.Marshal(game)
		if err != nil {
			log.Println(err)
		}
		email.SetBodyData(mail.TextPlain, info)
		//email.AddCc("another_you@example.com")
		email.SetSubject("Discount Notification")
		//email.SetBody(mail.TextHTML, htmlBody)
		// Send email
		err = email.Send(smtpClient)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("EMAIL NOTIFICATION SENT")
	}
	return err
}

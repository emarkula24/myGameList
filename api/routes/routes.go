package routes

import (
	"example.com/mygamelist/handler"
	"github.com/gorilla/mux"
)

func CreateGameSubrouter(router *mux.Router, game *handler.GameHandler) *mux.Router {
	s := router.PathPrefix("/games").Subrouter()
	s.HandleFunc("/search", game.Search).Methods("GET")
	s.HandleFunc("/game", game.SearchGame).Methods("GET")
	return s
}
func CreateUserSubrouter(router *mux.Router, user *handler.UserHandler) *mux.Router {
	s := router.PathPrefix("/user").Subrouter()
	s.HandleFunc("/register", user.Register).Methods("POST")
	s.HandleFunc("/login", user.Login).Methods("POST")
	s.HandleFunc("/refresh", user.Refresh).Methods("POST")
	return s

}

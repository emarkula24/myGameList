package routes

import (
	"net/http"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/middleware"
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
	s.HandleFunc("/logout", user.Logout).Methods("POST")
	s.HandleFunc("/users", user.GetUsers).Methods("GET")
	return s

}
func CreateListSubRouter(router *mux.Router, list *handler.ListHandler) *mux.Router {
	s := router.PathPrefix("/list").Subrouter()
	s.HandleFunc("", list.GetList).Methods("GET")
	s.HandleFunc("/game", list.GetListItem).Methods("GET")
	s.Handle("/add", middleware.VerifyJWTMiddleware(http.HandlerFunc(list.InsertToList))).Methods("POST")
	s.Handle("/update", middleware.VerifyJWTMiddleware(http.HandlerFunc(list.UpdateList))).Methods("PUT")
	s.Handle("/delete", middleware.VerifyJWTMiddleware(http.HandlerFunc(list.DeleteListItem))).Methods("DELETE")

	return s
}

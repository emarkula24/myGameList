package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/service"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const port string = ":8080"

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env variables")
	}

	frontUrl := os.Getenv("VITE_FRONTEND_URL")

	cfg := mysql.NewConfig()
	cfg.User = "mies"
	cfg.Passwd = "mies"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "test"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	client := &service.GiantBombClient{}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	h := handler.NewHandler(service)
	gameHandler := handler.NewGameHandler(client)
	searchSubRoute := router.PathPrefix("/games").Subrouter()
	searchSubRoute.HandleFunc("/", gameHandler.Search).Methods("GET")
	searchSubRoute.HandleFunc("/game", gameHandler.SearchGame).Methods("GET")

	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/refresh", h.Refresh).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(router)
	fmt.Printf("running server on port %s \n", port)
	http.ListenAndServe(port, corsHandler)
}

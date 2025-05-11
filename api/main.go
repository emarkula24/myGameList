package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	env := &utils.Env{DB: db}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.HandleFunc("/search", handler.Search).Methods("GET")
	router.HandleFunc("/register", handler.Register(env)).Methods("POST")
	router.HandleFunc("/login", handler.Login(env)).Methods("POST")
	fmt.Printf("running server on port %s \n", port)

	http.ListenAndServe(port, router)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

// type AppHandler func(http.ResponseWriter, *http.Request) error

// func WithErrorHandler(h AppHandler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if err := h(w, r); err != nil {
// 			log.Printf("handling %q: %v", r.RequestURI, err)
// 			errorutils.WriteJSONError(w, "something went wrong", http.StatusInternalServerError)
// 		}
// 	})
// }

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env variables")
	}

	frontUrl := os.Getenv("VITE_FRONTEND_URL")

	cfg := mysql.NewConfig()
	cfg.User = "mies"
	cfg.Passwd = "mies"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3308"
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
	service := service.NewUserService(repo)

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

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{frontUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	corsHandler := cors.Handler(router)

	srv := &http.Server{
		Addr: port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsHandler, // Pass our instance of gorilla/mux in.

	}
	fmt.Printf("running server on port %s \n", port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

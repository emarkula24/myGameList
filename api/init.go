package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/middleware"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/routes"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Defines application handlers.
type Handlers struct {
	user *handler.UserHandler
	game *handler.GameHandler
	list *handler.ListHandler
}

// SetUp initializes new application handler dependencies.
func SetUp(db *sql.DB) *Handlers {
	client := service.NewGiantBombClient()
	auth := &utils.AuthService{}
	repo := repository.NewRepository(db)
	listRepo := repository.NewListRepository(db)

	userService := service.NewUserService(repo, auth)
	listService := service.NewListService(listRepo, client)

	userHandler := handler.NewUserHandler(userService)
	gameHandler := handler.NewGameHandler(client)
	listHandler := handler.NewListHandler(listService)

	handlers := &Handlers{
		user: userHandler,
		game: gameHandler,
		list: listHandler,
	}

	return handlers

}

// NewDatabase creates a new database pool.
func NewDatabase() *sql.DB {

	u := os.Getenv("MYSQL_USER")
	p := os.Getenv("MYSQL_PASSWORD")
	a := os.Getenv("MYSQL_ADDRESS")
	n := os.Getenv("MYSQL_DATABASE")
	cfg := mysql.NewConfig()
	cfg.User = u
	cfg.Passwd = p
	cfg.Net = "tcp"
	cfg.Addr = a
	cfg.DBName = n

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

// Router creates a new router instance with subroutes.
func Router() *mux.Router {

	mode := os.Getenv("MODE")
	db := NewDatabase()
	handlers := SetUp(db)
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	if mode == "development" {

		router.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
			// Deletes all data in the development table if called in tests
			_, err := db.Exec("SET FOREIGN_KEY_CHECKS=0")
			if err != nil {
				log.Fatalf("Failed to disable FK checks: %v", err)
			}

			tables := []string{"user_games", "refreshtokens", "games", "users"}
			for _, t := range tables {
				_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", t))
				if err != nil {
					log.Fatalf("Failed to truncate %s: %v", t, err)
				}
			}

			_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
			if err != nil {
				log.Fatalf("Failed to enable FK checks: %v", err)
			}

			if err != nil {
				log.Printf("Failed to reset database: %v", err)
				w.WriteHeader(500)
				_, err = w.Write([]byte("Failed to reset database"))
				if err != nil {
					log.Fatalf("failed to resed database")
				}
				return
			}
			w.WriteHeader(200)
			_, err = w.Write([]byte("Database reset successfully"))
			if err != nil {
				log.Fatalf("failed to reset database")
			}
		}).Methods("POST")
	} else {
		log.Println("running on production mode, /reset not available")
	}

	routes.CreateGameSubrouter(router, handlers.game)
	routes.CreateUserSubrouter(router, handlers.user)
	routes.CreateListSubRouter(router, handlers.list)

	return router
}

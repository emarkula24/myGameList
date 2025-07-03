package main

import (
	"database/sql"
	"log"
	"os"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/routes"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Handlers struct {
	user *handler.UserHandler
	game *handler.GameHandler
}

func setUpDependencies(db *sql.DB) *Handlers {
	client := &service.GiantBombClient{}
	auth := &utils.AuthService{}
	repo2 := repository.NewRepository(db)
	service := service.NewUserService(repo2, auth)

	userHandler := handler.NewUserHandler(service)
	gameHandler := handler.NewGameHandler(client)

	handlers := &Handlers{
		user: userHandler,
		game: gameHandler,
	}

	return handlers

}

func setUpDatabase() *sql.DB {

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

func initializeServer() *mux.Router {
	mode := os.Getenv("MODE")

	// Development: load from file
	if mode != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Running in development mode, local .env file needed")
		}
	} else {
		// Optional: print loaded mode
		log.Println("Running in production mode.")
	}

	db := setUpDatabase()
	handlers := setUpDependencies(db)

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	routes.CreateGameSubrouter(router, handlers.game)
	routes.CreateUserSubrouter(router, handlers.user)

	return router
}

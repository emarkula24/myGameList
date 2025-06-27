package main

import (
	"database/sql"
	"log"

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

	return db
}

func initializeServer() *mux.Router {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env variables")
	}

	db := setUpDatabase()
	handlers := setUpDependencies(db)

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	routes.CreateGameSubrouter(router, handlers.game)
	routes.CreateUserSubrouter(router, handlers.user)

	return router
}

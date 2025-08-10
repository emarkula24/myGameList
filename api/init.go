package main

import (
	"database/sql"
	"log"
	"os"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/middleware"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/routes"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Defines application handlers.
type Handlers struct {
	user *handler.UserHandler
	game *handler.GameHandler
	list *handler.ListHandler
}

// SetUp initializes new application handler dependencies.
func SetUp(db *sql.DB) *Handlers {
	client := &service.GiantBombClient{}
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

	// Development: load from file
	if mode != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Running in development mode, local .env file needed")
		}
	} else {
		// Optional: print loaded mode
		log.Println("Running in production mode.")
	}

	db := NewDatabase()
	handlers := SetUp(db)

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)

	routes.CreateGameSubrouter(router, handlers.game)
	routes.CreateUserSubrouter(router, handlers.user)
	routes.CreateListSubRouter(router, handlers.list)

	return router
}

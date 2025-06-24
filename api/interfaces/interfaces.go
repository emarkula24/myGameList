package interfaces

import "net/http"

type Client interface {
	SearchGames(query string) (*http.Response, error)
}

type UserRepository interface {
	SelectUserByUsername(username string) (bool, error)
	InsertUser(username, email, hashedPassword string) (int64, error)
}
type UserService interface {
	RegisterUser(username, email, password string) (int64, error)
}
type Repository interface {
	UserRepository
	// add other entity interfaces here when needed
}

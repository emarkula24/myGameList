package interfaces

import "net/http"

type GiantBombClient interface {
	SearchGames(query string) (*http.Response, error)
}

// type UserRepository interface {
// 	SelectUserByUsername(username string) (bool, error)
// 	InsertUser(username, email, hashedPassword string) (int64, error)
// 	PasswordByUsername(username, password string) (string, error)
// 	SelectUserIdByUsername(username string) (int, error)
// 	InsertRefreshToken(userId int, refreshtoken, jti string) error
// 	RefreshTokenById(userId int) (string, string, error)
// }
// type UserService interface {
// 	RegisterUser(username, email, password string) (int64, error)
// 	LoginUser(username, password string) (string, error)
// 	StoreRefreshToken(username, refreshtoken, jti string) error
// 	FetchRefreshToken(username string) (string, error)
// }
// type Repository interface {
// 	UserRepository
// 	// add other entity interfaces here when needed
// }

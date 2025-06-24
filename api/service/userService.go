package service

import (
	"fmt"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	UserRepository *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{UserRepository: repo}
}
func (s *UserService) RegisterUser(username, email, password string) (int64, error) {

	isUser, err := s.UserRepository.SelectUserByUsername(username)
	if err != nil {
		return 0, fmt.Errorf("failed to select user: %w", err)
	}

	if isUser {
		return 0, errorutils.ErrUserExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userId, err := s.UserRepository.InsertUser(username, email, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userId, nil
}

func (s *UserService) LoginUser(username, password string) (string, error) {
	var secretKey = []byte("secret-key")
	expirationTime := time.Now().Add(5 * time.Minute).Unix()
	hashedPassword, err := s.UserRepository.PasswordByUsername(username, password)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve users password: %w", err)
	}
	if !utils.CheckPasswordHash(hashedPassword, password) {
		return "", errorutils.ErrPasswordMatch
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err

}

func (s *UserService) StoreRefreshToken(username, refreshToken, jti string) error {
	userId, err := s.UserRepository.SelectUserIdByUsername(username)
	if err != nil {
		return err
	}
	err = s.UserRepository.InsertRefreshToken(userId, refreshToken, jti)
	if err != nil {
		return fmt.Errorf("failed to insert refreshtoken: %w", err)
	}

	return nil
}
func (s *UserService) FetchRefreshToken(username string) (string, error) {
	userId, err := s.UserRepository.SelectUserIdByUsername(username)
	if err != nil {
		return "", fmt.Errorf("failed to fetch username: %w", err)
		// http.Error(w, "failed to retrieve userId", http.StatusUnauthorized)
	}
	_, jtiFromDb, err := s.UserRepository.RefreshTokenById(userId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch refreshtoken: %w", err)
	}
	return jtiFromDb, nil
}
